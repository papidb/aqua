package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/resources"
	"github.com/papidb/aqua/services/core/pkg/controllers"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
	"github.com/spf13/cobra"
)

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with cloud resources",
	Long:  fmt.Sprintf(`Generate and seed the database with %d randomly generated cloud resources for testing purposes.`, resources.DefaultMaxResources),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var env config.Env
		if err := config.LoadEnv(&env); err != nil {
			panic(err)
		}
		app, err := config.New(env)
		if err != nil {
			panic(err)
		}

		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}
		defer app.Database.DB.Close()

		// Retrieve the flag value for max resources
		maxResources, err := cmd.Flags().GetInt("max-resources")
		if err != nil {
			log.Fatalf("Error retrieving max-resources flag: %v", err)
		}

		err = resources.SeedResources(app, maxResources)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "Aqua CLI",
	Short: "Aqua CLI",
	Long: `Aqua CLI is a command-line interface for managing cloud resources.
	You can use it to seed the database with sample cloud resources.
	Start by running "aqua seed" to generate and seed the database with 100 cloud resources.
	You can also use it to start the server by running "aqua server".
	`,
}

func main() {
	// Add commands to rootCmd
	SeedCmd.Flags().Int("max-resources", resources.DefaultMaxResources, "Number of cloud resources to generate (default: 100 if not provided)")
	rootCmd.AddCommand(SeedCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Start the Gin server",
		Run:   startServer,
	})

	// Execute Cobra
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer(cmd *cobra.Command, args []string) {
	var err error
	var env config.Env
	if err := config.LoadEnv(&env); err != nil {
		panic(err)
	}
	app, err := config.New(env)
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	middlewares.PrepareRequest(app, router, app.Logger)
	handler := controllers.MountRoutes(app, router)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Env.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	defer func() {
		if err := app.Database.Close(); err != nil {
			app.Logger.Err(err).Msg("failed to disconnect from postgres cleanly")
		}
	}()
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
