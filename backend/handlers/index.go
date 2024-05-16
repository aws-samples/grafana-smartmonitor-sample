package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"genai.ai/automonitor/service"
)

// Index is a handler function that returns a welcome message
// for the Grafana Screenshot API.
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Grafana Screenshot API",
	})
}

// Config is a handler function that returns the global configuration
// for the Grafana Screenshot API.
func Config(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"config": service.GetGlobalConfig(),
	})
}
