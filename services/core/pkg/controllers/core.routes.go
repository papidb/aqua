package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
)

func MountRoutes(app *config.App, r *gin.Engine) http.Handler {
	gin.ForceConsoleColor()

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
