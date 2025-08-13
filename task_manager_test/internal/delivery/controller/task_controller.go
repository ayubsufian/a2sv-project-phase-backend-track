package controller

import (
	"errors"
	"net/http"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

// TaskController wraps use case interfaces for task operations.
type TaskController struct {
	taskUC usecase.TaskUsecase
}

// NewTaskController creates a new Handler given Task use cases.
func NewTaskController(t usecase.TaskUsecase) *TaskController {
	return &TaskController{taskUC: t}
}

// TaskResponse defines the JSON structure for task data returned in API responses.
type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"duedate"`
	Status      string    `json:"status"`
}

// mapToTaskResponse converts a domain.Task into a TaskResponse for API output.
func mapToTaskResponse(t domain.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Status:      t.Status,
	}
}

// GetTasks retrieves all tasks via taskUC.List and returns them as JSON.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.taskUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve tasks"})
		return
	}

	// Map tasks to TaskResponse
	responses := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = mapToTaskResponse(t)
	}

	c.JSON(http.StatusOK, responses)
}

// GetTask retrieves a single task by ID via taskUC.Get.
func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskUC.Get(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID format"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve task"})
		}
		return
	}
	c.JSON(http.StatusOK, mapToTaskResponse(task))
}

// CreateTask handles the creation of a new task.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var body struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		DueDate     *time.Time `json:"duedate" binding:"required"`
		Status      string     `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := domain.Task{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     *body.DueDate,
		Status:      body.Status,
	}
	created, err := tc.taskUC.Create(c.Request.Context(), task)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrTaskAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "a task with these details already exists"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, mapToTaskResponse(created))
}

// UpdateTask updates an existing task identified by URL param ID.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		DueDate     *time.Time `json:"duedate" binding:"required"`
		Status      string     `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := domain.Task{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
		DueDate:     *body.DueDate,
		Status:      body.Status,
	}
	updated, err := tc.taskUC.Update(c.Request.Context(), task)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID format"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
		}
		return
	}
	c.JSON(http.StatusOK, mapToTaskResponse(updated))
}

// DeleteTask deletes a task by ID via taskUC.Delete.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := tc.taskUC.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, usecase.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID format"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// AdminDashboard handles admin-only access.
func (tc *TaskController) AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin"})
}
