package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/customers"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
)

func MountRoutes(app *config.App, r *gin.Engine) http.Handler {
	gin.ForceConsoleColor()

	customerRepo := customers.NewRepo(app.Database.DB)
	customerService := customers.NewService(app.Database.DB, customerRepo)

	r.GET("/", helloWorldHandler)
	r.GET("/health", func(ctx *gin.Context) {
		healthHandler(ctx, app)
	})

	// create customer
	r.POST(
		"/customers",
		middlewares.ValidationMiddleware(&customers.CreateCustomerDTO{}),
		createCustomerHandler(app, customerService),
	)
	// add cloud resource to customer
	r.POST("/customers/:customer_id/resources", addCloudResourceHandler)
	// Fetch Cloud Resources by Customer
	r.GET("/customers/:customer_id/resources", fetchCloudResourcesHandler)
	// Update Resource Information
	r.PUT("/resources/:resource_id", updateResourceHandler)
	// Delete a Resource
	r.DELETE("/resources/:resource_id", deleteResourceHandler)

	return r
}
