package controllers

import (
	"net/http"
	"task_manager_clean/domain"
	"task_manager_clean/usecases"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler wraps use case interfaces for user and task operations.
type Handler struct {
	userUC usecases.UserUsecase
	taskUC usecases.TaskUsecase
}

// NewHandler creates a new Handler given User and Task use cases.
func NewHandler(u usecases.UserUsecase, t usecases.TaskUsecase) *Handler {
	return &Handler{userUC: u, taskUC: t}
}

// Register handles new user registration requests.
func (h *Handler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := domain.User{
		Username: body.Username,
		Password: body.Password,
		Role:     body.Role,
	}
	if err := h.userUC.Register(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User Registered successfully"})
}

// Login handles user authentication.
func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.userUC.Login(c.Request.Context(), body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
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
func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.taskUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.taskUC.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, mapToTaskResponse(task))
}

// CreateTask handles the creation of a new task.
func (h *Handler) CreateTask(c *gin.Context) {
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
	created, err := h.taskUC.Create(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapToTaskResponse(created))
}

// UpdateTask updates an existing task identified by URL param ID.
func (h *Handler) UpdateTask(c *gin.Context) {
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
	updated, err := h.taskUC.Update(c.Request.Context(), task)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, mapToTaskResponse(updated))
}

// DeleteTask deletes a task by ID via taskUC.Delete.
func (h *Handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.taskUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// AdminDashboard handles admin-only access.
func (h *Handler) AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin"})
}
