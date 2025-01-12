package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/customers"
	"github.com/papidb/aqua/pkg/entities/resources"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
)

func MountRoutes(app *config.App, r *gin.Engine) http.Handler {
	gin.ForceConsoleColor()

	customerRepo := customers.NewRepo(app.Database.DB)
	resourceRepo := resources.NewRepo(app.Database.DB)

	customerService := customers.NewService(app.Database.DB, customerRepo, resourceRepo)
	resourceService := resources.NewService(app.Database.DB, resourceRepo)

	r.GET("/", helloWorldHandler)
	r.GET("/health", func(ctx *gin.Context) {
		healthHandler(ctx, app)
	})

	// create customer
	r.POST(
		"/customers",
		middlewares.ValidationBodyMiddleware(&customers.CreateCustomerDTO{}),
		createCustomerHandler(app, customerService),
	)
	// add cloud resource to customer
	r.POST(
		"/customers/:customer_id/resources",
		// TODO: VALIDATE CUSTOMER ID
		middlewares.ValidationBodyMiddleware(&customers.AddResourceToCustomerDTO{}),
		addCloudResourceHandler(app, customerService),
	)
	// Fetch Cloud Resources by Customer
	r.GET("/customers/:customer_id/resources", fetchCloudResourcesHandler(app, customerService))
	// Update Resource Information
	r.PUT("/resources/:resource_id",
		middlewares.ValidationBodyMiddleware(&resources.UpdateResourceDTO{}),
		updateResourceHandler(app, resourceService))
	// Delete a Resource
	r.DELETE("/resources/:resource_id", deleteResourceHandler(app, resourceService, customerService))

	return r
}
