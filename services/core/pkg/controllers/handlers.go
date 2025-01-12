package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/api"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/customers"
	"github.com/papidb/aqua/pkg/entities/resources"
)

func helloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

var errorMapping = map[error]int{
	customers.ErrExistingEmailOrName{}:      http.StatusBadRequest,
	customers.ErrExistingCustomerResource{}: http.StatusBadRequest,
	api.ErrCustomerNotFound:                 http.StatusNotFound,
	api.ErrResourceNotFound:                 http.StatusNotFound,
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

		if api.HandleMappedErrors(c, err, errorMapping) {
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
		if api.HandleMappedErrors(c, err, errorMapping) {
			return
		}

		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not add your resource to customer.",
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

func fetchCloudResourcesHandler(_ *config.App, customerService *customers.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.Param("customer_id")

		resources, err := customerService.FetchCloudResourcesByCustomerID(c.Request.Context(), customerID)
		if api.HandleMappedErrors(c, err, errorMapping) {
			return
		}

		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not fetch the resources.",
				Err:     err,
			})
			return
		}

		api.Success(c.Request, c.Writer, &api.AppResponse{
			Message: "Resources fetched successfully",
			Data:    resources,
			Code:    http.StatusCreated,
		})
	}
}

func updateResourceHandler(_ *config.App, resourceService *resources.ResourceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto resources.UpdateResourceDTO
		c.ShouldBindJSON(&dto)
		resource_id := c.Param("resource_id")

		resource, err := resourceService.UpdateResource(c.Request.Context(), resource_id, dto)
		if api.HandleMappedErrors(c, err, errorMapping) {
			return
		}

		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not update resource.",
				Err:     err,
			})
			return
		}

		api.Success(c.Request, c.Writer, &api.AppResponse{
			Message: "Updated resource successfully",
			Data:    resource,
			Code:    http.StatusCreated,
		})
	}
}

func deleteResourceHandler(_ *config.App, resourceService *resources.ResourceService, customerService *customers.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource_id := c.Param("resource_id")

		ctx := c.Request.Context()
		resource, err := resourceService.DeleteResource(ctx, resource_id)
		if api.HandleMappedErrors(c, err, errorMapping) {
			return
		}

		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not delete resource.",
				Err:     err,
			})
			return
		}

		err = customerService.DeleteCustomerResource(ctx, resource_id)
		if api.HandleMappedErrors(c, err, errorMapping) {
			return
		}

		if err != nil {
			api.Error(c.Request, c.Writer, api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not delete resource.",
				Err:     err,
			})
			return
		}

		api.Success(c.Request, c.Writer, &api.AppResponse{
			Message: "Deleted resource successfully",
			Data:    resource,
			Code:    http.StatusCreated,
		})
	}
}
