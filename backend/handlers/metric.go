package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"genai.ai/automonitor/service"
)

// RunMetricByID runs a monitor metric by its ID and updates the metric status and screenshot.
func RunMetricByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	exitsMetric, err := service.GetMonitorMetric(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor Metric not found"})
	}

	
	exitsConnection, err := service.GetMonitorConnectionByName(exitsMetric.ConnectionName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Connection %s not found",exitsMetric.ConnectionName)})
	}

	
	fmt.Println(exitsConnection)


	result, err := service.RunCondition(&exitsMetric,&exitsConnection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ScreentCapture Faild %s", err)})
		return
	}

	var checkResult service.CheckResult

	err = yaml.Unmarshal([]byte(result[1]), &checkResult)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"checkResult Unmarshal error": err.Error()})
		return
	}

	exitsMetric.Status = checkResult.Result.Pass
	exitsMetric.StatusDesc = checkResult.Result.Reason
	exitsMetric.CheckDate = time.Now()
	exitsMetric.Screen = result[0]

	err = service.UpdateMonitorMetric(exitsMetric)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exitsMetric)

}

// GetMonitorMetricByID retrieves a monitor metric by its ID.
func GetMonitorMetricByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	exitsMetric, err := service.GetMonitorMetric(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor Metric not found"})
	}

	c.JSON(http.StatusCreated, exitsMetric)

}

// CreateMonitorMetric creates a new monitor metric.
func CreateMonitorMetric(c *gin.Context) {
	var newMonitorMetric service.MonitorMetric
	if err := c.BindJSON(&newMonitorMetric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connection,err:=service.GetMonitorConnectionByName(newMonitorMetric.ConnectionName)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return 
	}
	// Add the new monitor Metric to the slice
	newMonitorMetric.DashboardURL=service.ReplaceHost(connection.URL,newMonitorMetric.DashboardURL)

	id, err := service.InsertMonitorMetric(newMonitorMetric)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	newMonitorMetric.ID = id

	metric, _ := service.GetMonitorMetric(id)
	c.JSON(http.StatusCreated, metric)
}

// UpdateMonitorMetric updates an existing monitor metric.
func UpdateMonitorMetric(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var newMonitorMetric service.MonitorMetric
	if err := c.BindJSON(&newMonitorMetric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connection,err:=service.GetMonitorConnectionByName(newMonitorMetric.ConnectionName)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return 
	}
	// Add the new monitor Metric to the slice
	newMonitorMetric.DashboardURL=service.ReplaceHost(connection.URL,newMonitorMetric.DashboardURL)


	err = service.UpdateMonitorMetric(newMonitorMetric)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor Metric not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d, update successully", id)})
}

// DeleteMonitorMetric deletes a monitor metric by its ID.
func DeleteMonitorMetric(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = service.DeleteMonitorMetric(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d, deleted successully", id)})
}

// GetMetrics retrieves monitor metrics based on the provided page number and project.
func GetMetrics(c *gin.Context) {

	var metrics []service.MonitorMetric
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	project := c.Query("project")

	if project == "" {
		metrics, err = service.GetMetrics(int(page))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		metrics, err = service.GetMetricsByProject(project)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

// GetProjects retrieves distinct project names from the monitor_metrics table.
func GetProjects(c *gin.Context) {

	var metrics []map[string]string

	metrics, err := service.GetDistinct("project", "monitor_metrics")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": metrics})
}

// GetCatalogs retrieves distinct catalog names from the monitor_metrics table.
func GetCatalogs(c *gin.Context) {

	var items []map[string]string

	items, err := service.GetDistinct("catalog", "monitor_metrics")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
