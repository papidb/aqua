package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Env is the expected config values from the process's environment
type ApplicationEnvironment string

const (
	Dev        = "dev"
	Staging    = "staging"
	Production = "production"
)

type Env struct {
	AppEnv ApplicationEnvironment `default:"dev" split_words:"true"`
	Name   string                 `envconfig:"SERVICE_NAME" required:"true"`
	Port   int                    `required:"true"`

	PostgresHost       string `required:"true" split_words:"true"`
	PostgresPort       string `required:"true" split_words:"true"`
	PostgresPoolSize   int    `required:"true" split_words:"true"`
	PostgresSecureMode bool   `required:"true" split_words:"true"`
	PostgresUser       string `required:"true" split_words:"true"`
	PostgresPassword   string `required:"true" split_words:"true"`
	PostgresDatabase   string `required:"true" split_words:"true"`
	PostgresDebug      bool   `default:"false" split_words:"true"`

	RabbitMQHost     string `envconfig:"RABBIT_MQ_HOST" required:"true" split_words:"true"`
	RabbitMQPort     string `envconfig:"RABBIT_MQ_PORT" required:"true" split_words:"true"`
	RabbitMQUser     string `envconfig:"RABBITMQ_DEFAULT_USER" required:"true" split_words:"true"`
	RabbitMQPassword string `envconfig:"RABBITMQ_DEFAULT_PASS" required:"true" split_words:"true"`
	RabbitMQQueue    string `envconfig:"RABBIT_MQ_QUEUE" required:"true" split_words:"true"`

	RPCPort string `envconfig:"RPC_PORT" required:"true" split_words:"true"`
}

// LoadEnv loads environment variables into env
func LoadEnv(env *Env) error {
	// try to load from .env first
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	err := godotenv.Load(".env." + appEnv + ".local")
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok || !errors.Is(perr.Unwrap(), os.ErrNotExist) {
			return err
		}
	}

	if appEnv != "test" {
		godotenv.Load(".env.local")
	}
	err = godotenv.Load(".env." + appEnv)
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok || !errors.Is(perr.Unwrap(), os.ErrNotExist) {
			return err
		}
	}
	err = godotenv.Load() // The Original .env
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok || !errors.Is(perr.Unwrap(), os.ErrNotExist) {
			return err
		}
	}

	if err := envconfig.Process("", env); err != nil {
		return err
	}

	return nil
}
