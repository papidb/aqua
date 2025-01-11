package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer(env config.Env) *http.Server {
	NewServer := &Server{
		port: env.Port,
		db:   database.New(env),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
