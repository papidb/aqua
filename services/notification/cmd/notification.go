package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/notification"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
	Server "github.com/papidb/aqua/services/core/pkg/server"
	"github.com/papidb/aqua/services/notification/pkg/controllers"
	"google.golang.org/grpc"
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
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Env.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server
	go func() {
		log.Printf("HTTP server is running on port %d", app.Env.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start the gRPC server
	go func() {
		rpcPort := fmt.Sprintf(":%s", app.Env.RPCPort)
		listener, err := net.Listen("tcp", rpcPort)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		notificationService := &notification.NotificationServiceImpl{
			NotificationService: service,
		}
		// Create a NotificationServiceServer implementation
		notification.RegisterNotificationServiceServer(grpcServer, notificationService)

		// // Enable reflection for debugging (optional)
		// reflection.Register(grpcServer)

		log.Printf("gRPC server is running on port %s", rpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Channel to signal shutdown completion
	done := make(chan bool, 1)

	// Start RabbitMQ listener in a separate goroutine
	go func() {
		config.ListenForNotifications(env, func(msg string) {
			app.Logger.Printf("Received message: %s", msg)
			// Add to notification service
			service.AddNotification("user", msg)
		})
	}()

	// Run graceful shutdown in a separate goroutine
	go func() {
		Server.GracefulShutdown(httpServer, done)
	}()

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete")
}
