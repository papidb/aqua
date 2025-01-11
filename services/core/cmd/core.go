package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/services/core/server"
	"github.com/uptrace/bun"
)

func main() {
	var err error

	var env config.Env
	if err = config.LoadEnv(&env); err != nil {
		panic(err)
	}

	log := config.NewLogger(env.Name)

	// _ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// connect to postgresql
	var db *bun.DB
	if db, err = config.SetupDB(env); err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Err(err).Msg("failed to disconnect from postgres cleanly")
		}
	}()
	log.Info().Msg("successfully connected to postgres and has run migrations")

	// app := &config.App{
	// 	Env: &env, DB: db, Logger: log,
	// }

	server := server.NewServer(env)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
