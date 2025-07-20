package controllers

import (
	"net/http"
	"strings"

	"task_manager_mongodb/data"
	"task_manager_mongodb/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// parseObjectID extracts and validates ObjectID from URL params
func parseObjectID(c *gin.Context) (primitive.ObjectID, bool) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return primitive.NilObjectID, false
	}
	return objID, true
}

// handleValidationError handles request validation errors.
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

// Retrieves all tasks from the database.
func GetTasks(c *gin.Context) {
	tasks, err := data.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Retrieves a single task by its ID.
func GetTask(c *gin.Context) {
	objID, ok := parseObjectID(c)
	if !ok {
		return
	}

	task, err := data.GetByID(objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Binds JSON input, validates it, and inserts the task into the database.
func CreateTask(c *gin.Context) {
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		handleValidationError(c, err)
		return
	}

	created, err := data.Create(t)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "task already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

// Updates a task by its ID after validating the input.
func UpdateTask(c *gin.Context) {
	objID, ok := parseObjectID(c)
	if !ok {
		return
	}

	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		handleValidationError(c, err)
		return
	}

	updated, err := data.Update(objID, t)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate value for unique field"})
		} else if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Deletes a task by its ID.
func DeleteTask(c *gin.Context) {
	objID, ok := parseObjectID(c)
	if !ok {
		return
	}

	if err := data.Delete(objID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
