package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/papidb/aqua/pkg/config"
)

type Server struct {
	App *config.App
}

func NewServer(app *config.App) *http.Server {
	NewServer := &Server{
		App: app,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Env.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
