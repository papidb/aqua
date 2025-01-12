package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/cli"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/resources"
	"github.com/papidb/aqua/services/core/pkg/controllers"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
	Server "github.com/papidb/aqua/services/core/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	// Add commands to rootCmd
	cli.SeedCmd.Flags().Int("max-resources", resources.DefaultMaxResources, "Number of cloud resources to generate (default: 100 if not provided)")
	cli.RootCmd.AddCommand(cli.SeedCmd)
	cli.RootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Start the Gin server",
		Run:   startServer,
	})

	// Execute Cobra
	if err := cli.RootCmd.Execute(); err != nil {
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
	go Server.GracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
