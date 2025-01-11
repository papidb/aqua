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

	// create customer
	r.POST("/customers", dummyHandler)
	// add cloud resource to customer
	r.POST("/customers/:customer_id/resources", dummyHandler)
	// Fetch Cloud Resources by Customer
	r.GET("/customers/:customer_id/resources", dummyHandler)
	// Update Resource Information
	r.PUT("/customers/:customer_id/resources/:resource_id", dummyHandler)
	// Delete a Resource
	r.DELETE("/customers/:customer_id/resources/:resource_id", dummyHandler)

	return r
}
