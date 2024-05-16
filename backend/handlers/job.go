package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"genai.ai/automonitor/service"
	"github.com/gin-gonic/gin"
)



// GetMonitorJob retrieves a monitor job by its ID.
func GetMonitorJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	job, err := service.GetMonitorJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// GetJobStatus retrieves the status of a monitor job by its project name.
func GetJobStatus(c *gin.Context) {
	projectName := c.Param("project")

	job, ok := service.GetJobByProject(projectName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": *job})
}

// CreateMonitorJob creates a new monitor job.
func CreateMonitorJob(c *gin.Context) {
	var newJob service.MonitorJob
	if err := c.BindJSON(&newJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := service.InsertMonitorJob(newJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newJob.ID = id
	fmt.Printf("newJob.Enable: %v\n", newJob.Enable)
	if newJob.Enable {
		service.AddJob(newJob.Project, newJob.Cron)
	}
	c.JSON(http.StatusCreated, newJob)
}

// UpdateMonitorJob updates an existing monitor job.
func UpdateMonitorJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedJob service.MonitorJob
	if err := c.BindJSON(&updatedJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedJob.ID = id
	err = service.UpdateMonitorJob(updatedJob)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Monitor job not found"})
		return
	}

	if updatedJob.Enable {
		service.AddJob(updatedJob.Project, updatedJob.Cron)
	} else {
		service.RemoveJobByProject(updatedJob.Project)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monitor job updated successfully"})
}

// DeleteMonitorJob deletes a monitor job by its ID.
func DeleteMonitorJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = service.DeleteMonitorJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monitor job deleted successfully"})
}

// GetMonitorJobs retrieves all monitor jobs.
func GetMonitorJobs(c *gin.Context) {
	jobs, err := service.GetMonitorJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
