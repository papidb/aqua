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

func createCustomerHandler(_ *config.App, customerService *customers.CustomerService) gin.HandlerFunc {
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

func addCloudResourceHandler(_ *config.App, customerService *customers.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto customers.AddResourceToCustomerDTO
		c.ShouldBindJSON(&dto)
		customerID := c.Param("customer_id")

		customer, resource, err := customerService.AddResourceToCustomer(c.Request.Context(), customerID, dto)

		if errors.Is(err, customers.ErrExistingCustomerResource{}) {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Err:     err,
			})
			return
		}

		if (err != nil) || (customer == nil) {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not add your resource to your customer.",
				Err:     err,
			})
			return
		}

		data := make(map[string]interface{})
		data["customer"] = customer
		data["resource"] = resource

		api.Success(c.Request, c.Writer, &api.AppResponse{
			Message: "Resource added to customer successfully",
			Data:    data,
			Code:    http.StatusCreated,
		})
	}
}

func fetchCloudResourcesHandler(_ *config.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "dummy",
		})
	}
}

func updateResourceHandler(_ *config.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "dummy",
		})
	}
}

func deleteResourceHandler(_ *config.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "dummy",
		})
	}
}
