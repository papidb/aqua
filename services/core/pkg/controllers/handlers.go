package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

func createCustomerHandler(c *gin.Context) {
	var req customers.CreateCustomerDTO
	_ = c.ShouldBindJSON(&req) // Already validated by middleware
	fmt.Println(&req)
	// c.JSON(http.StatusOK, gin.H{"message": "Customer created successfully", "customer": req})
	c.JSON(http.StatusOK, gin.H{"message": "Validation passed"})
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
