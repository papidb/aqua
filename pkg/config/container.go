package config

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// var
func MustStartPostgresContainer(env *Env) (func(ctx context.Context, opts ...testcontainers.TerminateOption) error, error) {

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
