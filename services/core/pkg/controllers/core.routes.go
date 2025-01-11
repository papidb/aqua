package controllers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
)

func MountRoutes(app *config.App) http.Handler {
	gin.ForceConsoleColor()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", helloWorldHandler)

	r.GET("/health", func(ctx *gin.Context) {
		healthHandler(ctx, app)
	})

	return r
}

func helloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func healthHandler(c *gin.Context, app *config.App) {
	c.JSON(http.StatusOK, app.Database.Health())
}
