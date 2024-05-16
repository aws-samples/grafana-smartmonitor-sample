package handlers

import (
	"net/http"
	"strconv"

	"genai.ai/automonitor/service"
	"github.com/gin-gonic/gin"
)

// GetMonitorConnection retrieves a monitor connection by ID
func GetMonitorConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	conn, err := service.GetMonitorConnection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor connection not found"})
		return
	}

	c.JSON(http.StatusOK, conn)
}

// CreateMonitorConnection creates a new monitor connection
func CreateMonitorConnection(c *gin.Context) {
	var newConn service.MonitorConnection
	if err := c.BindJSON(&newConn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := service.InsertMonitorConnection(newConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newConn.ID = id
	c.JSON(http.StatusCreated, newConn)
}

// UpdateMonitorConnection updates an existing monitor connection
func UpdateMonitorConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedConn service.MonitorConnection
	if err := c.BindJSON(&updatedConn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedConn.ID = id
	err = service.UpdateMonitorConnection(updatedConn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor connection not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monitor connection updated successfully"})
}

// DeleteMonitorConnection deletes a monitor connection by ID
func DeleteMonitorConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = service.DeleteMonitorConnection(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monitor connection deleted successfully"})
}


func GetMonitorConnectionName(c *gin.Context) {

	var metrics []map[string]string

	metrics, err := service.GetDistinct("conn_name", "monitor_connections")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connections": metrics})
}

// GetMonitorConnections retrieves all monitor connections
func GetMonitorConnections(c *gin.Context) {
	conns, err := service.GetMonitorConnections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connections": conns})
}
