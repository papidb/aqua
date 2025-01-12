package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/api"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/customers"
)

func helloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func healthHandler(c *gin.Context, app *config.App) {
	c.JSON(http.StatusOK, app.Database.Health())
}

func createCustomerHandler(app *config.App, customerService *customers.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req customers.CreateCustomerDTO
		c.ShouldBindJSON(&req)

		ctx := c.Request.Context()
		customer, err := customerService.CreateCustomer(ctx, req)
		if errors.Is(err, customers.ErrExistingEmailOrName{}) {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Err:     err,
			})
			return
		}
		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not create your customer.",
				Err:     err,
			})
			return
		}

		api.Success(c.Request, c.Writer, &api.AppResponse{
			Message: "Customer created successfully",
			Data:    customer,
			Code:    http.StatusCreated,
		})
	}

}

func addCloudResourceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "dummy",
	})
}

func fetchCloudResourcesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "dummy",
	})
}

func updateResourceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "dummy",
	})
}

func deleteResourceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "dummy",
	})
}
