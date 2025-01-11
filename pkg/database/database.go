package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"github.com/papidb/aqua/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/papidb/aqua/pkg/internal"
)

// migrateDB syncs the `dbname` with the migrations specified in cwd()/sql
// directory. It retruns all possible errors due to driver setup, migrator setup
// as well as possible errors while running the migrations (e.g syntax errors)
func migrateDB(dbname string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	dir := internal.GetPackagePath()
	dir = filepath.Join(dir, "sql/migrations")
	migrationsDir := fmt.Sprintf("file:///%s", dir)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsDir,
		dbname, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *bun.DB
}

var (
	database   string
	dbInstance *service
)

func New(env config.Env) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	database = env.PostgresDatabase
	connStr := dsnFromEnv(env)

	// attempts to create db pool
	sqldb, err := newStdDB(connStr, env)
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, pgdialect.New())

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	if err := migrateDB(env.PostgresDatabase, sqldb); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func dsnFromEnv(env config.Env) string {

	_, err := strconv.Atoi(env.PostgresPort)
	if err != nil {
		panic(fmt.Sprintf("Invalid Postgres port: %s", env.PostgresPort))
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.PostgresUser, env.PostgresPassword, env.PostgresHost, env.PostgresPort, env.PostgresDatabase)
}

// newStdDB creates a new database connection pool and returns a new sql.DB
func newStdDB(connStr string, env config.Env) (*sql.DB, error) {
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	// create std db from pgxpool since bun only works with database/sql instance
	sqldb := stdlib.OpenDBFromPool(dbpool, stdlib.OptionBeforeConnect(func(ctx context.Context, cc *pgx.ConnConfig) error {
		if !env.PostgresSecureMode {
			cc.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
		return nil
	}))

	return sqldb, nil
}
