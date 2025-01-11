package config

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type App struct {
	Env    *Env
	DB     *bun.DB
	wg     sync.WaitGroup
	Logger zerolog.Logger
}

func HealthChecker(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if err := app.DB.Ping(); err != nil {
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
