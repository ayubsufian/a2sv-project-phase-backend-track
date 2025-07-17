package controllers

import (
	"net/http"
	"strconv"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// handleValidationError parses and returns validation errors
func handleValidationError(c *gin.Context, err error) {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errs := make(map[string]string)
		for _, fe := range ve {
			errs[fe.Field()] = fe.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// GetTasks returns all tasks
func GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, data.GetAll())
}

// GetTask returns a task by its ID
func GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}
	task, err := data.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles task creation
func CreateTask(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		handleValidationError(c, err)
		return
	}
	created := data.Create(t)
	c.JSON(http.StatusCreated, created)
}

// UpdateTask modifies an existing task
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		handleValidationError(c, err)
		return
	}
	updated, err := data.Update(id, t)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteTask removes a task by ID
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}
	if err := data.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
