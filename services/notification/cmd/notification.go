package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
	Server "github.com/papidb/aqua/services/core/pkg/server"
	"github.com/papidb/aqua/services/notification/pkg/controllers"
)

func main() {
	// Load environment variables
	var env config.Env
	if err := config.LoadEnv(&env); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Initialize the application
	app, err := config.New(env)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer func() {
		if err := app.Database.Close(); err != nil {
			app.Logger.Err(err).Msg("Failed to disconnect from PostgreSQL cleanly")
		}
	}()

	// Set up the Gin router
	router := gin.Default()
	middlewares.PrepareRequest(app, router, app.Logger)
	handler, service := controllers.MountRoutes(app, router)

	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Env.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Channel to signal shutdown completion
	done := make(chan bool, 1)

	// Start RabbitMQ listener in a separate goroutine
	go func() {
		config.ListenForNotifications(env, func(msg string) {
			app.Logger.Printf("Received message: %s", msg)
			// add to notification service
			service.AddNotification("user", msg)
		})
	}()

	// Run graceful shutdown in a separate goroutine
	go func() {
		Server.GracefulShutdown(server, done)
	}()

	// Start the HTTP server
	log.Printf("Server is running on port %d", app.Env.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete")
}
