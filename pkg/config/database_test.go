package config

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

var env = &Env{
	PostgresDatabase: "database",
	PostgresPassword: "password",
	PostgresUser:     "user",
}

func TestMain(m *testing.M) {

	teardown, err := MustStartPostgresContainer(env)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	log.Println("started postgres container")

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
	srv := NewDB(*env)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := NewDB(*env)
	stats := srv.Health()
	log.Println("started postgres container")

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
	srv := NewDB(*env)
	log.Println("started postgres container")

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
