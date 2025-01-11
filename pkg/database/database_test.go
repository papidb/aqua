package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/papidb/aqua/pkg/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var env = config.Env{
	PostgresDatabase: "database",
	PostgresPassword: "password",
	PostgresUser:     "user",
}

func mustStartPostgresContainer() (func(ctx context.Context, opts ...testcontainers.TerminateOption) error, error) {

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(env.PostgresDatabase),
		postgres.WithUsername(env.PostgresUser),
		postgres.WithPassword(env.PostgresPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	database = env.PostgresDatabase

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	env.PostgresHost = dbHost
	env.PostgresPort = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	ctx := context.Background()
	if teardown != nil && teardown(ctx,
		testcontainers.RemoveVolumes(),
		testcontainers.StopContext(ctx),
	) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New(env)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New(env)
	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New(env)

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
