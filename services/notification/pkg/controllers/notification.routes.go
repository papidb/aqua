package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/notification"
)

func MountRoutes(app *config.App, router *gin.Engine) (http.Handler, *notification.NotificationService) {
	gin.ForceConsoleColor()
	notificationService := notification.NewNotificationService()

	router.POST("/notifications/:user_id", func(c *gin.Context) {
		user_id := c.Param("user_id")
		var body struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		notificationService.AddNotification(user_id, body.Message)
		c.JSON(http.StatusCreated, gin.H{"message": "Notification added"})
	})

	router.GET("/notifications/:user_id", func(c *gin.Context) {
		user_id := c.Param("user_id")
		notifications := notificationService.GetNotifications(user_id)
		c.JSON(http.StatusOK, notifications)
	})

	router.DELETE("/notifications/:user_id/:notification_id", func(c *gin.Context) {
		user_id := c.Param("user_id")
		notification_id := c.Param("notification_id")
		notificationService.ClearNotification(user_id, notification_id)
		c.JSON(http.StatusOK, gin.H{"message": "Notification cleared"})
	})

	router.DELETE("/notifications/:user_id", func(c *gin.Context) {
		user_id := c.Param("user_id")
		notificationService.ClearAllNotifications(user_id)
		c.JSON(http.StatusOK, gin.H{"message": "All notifications cleared"})
	})

	return router, notificationService
}
