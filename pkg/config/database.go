package config

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/papidb/aqua/pkg/internal"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
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

func SetupDB(env Env) (*bun.DB, error) {
	connStr := dsnFromEnv(env)

	// attempts to create db pool
	sqldb, err := newStdDB(connStr, env)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if env.PostgresDebug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	// auto migrate db on startuo
	if err := migrateDB(env.PostgresDatabase, sqldb); err != nil {
		return db, err
	}

	return db, nil
}

func dsnFromEnv(env Env) string {
	sslMode := "allow"
	if env.PostgresSecureMode {
		sslMode = "require"
	}

	port, err := strconv.Atoi(env.PostgresPort)
	if err != nil {
		panic(fmt.Sprintf("Invalid Postgres port: %s", env.PostgresPort))
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?application_name=%s&sslmode=%s&pool_max_conns=%d",
		env.PostgresUser, env.PostgresPassword, env.PostgresHost,
		port, env.PostgresDatabase, env.Name, sslMode, env.PostgresPoolSize,
	)
}

// newStdDB creates a new database connection pool and returns a new sql.DB
func newStdDB(connStr string, env Env) (*sql.DB, error) {
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
