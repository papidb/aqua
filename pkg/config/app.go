package config

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

type App struct {
	Env      *Env
	Database *Service
	wg       sync.WaitGroup
	Logger   zerolog.Logger
}

func New(env Env) (*App, error) {

	log := NewLogger(env.Name)

	// Connect to the database
	db := NewDB(env)

	log.Info().Msg("successfully connected to postgres and has run migrations")

	return &App{
		Env:      &env,
		Database: db,
		Logger:   log,
	}, nil
}

func HealthChecker(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if err := app.Database.DB.Ping(); err != nil {
			http.Error(w, "Could not reach database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Up and Running!"))
	}
}

func (app *App) Background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		// recover any panic
		defer func() {
			if err := recover(); err != nil {
				app.Logger.Err(fmt.Errorf("%s", err)).Msg("Failed to run background job")
			}
		}()

		// run the function
		fn()
	}()
}
