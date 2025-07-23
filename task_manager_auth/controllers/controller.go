package controllers

import (
	"net/http"
	"strings"
	"task_manager_auth/data"
	"task_manager_auth/middleware"
	"task_manager_auth/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// parseObjectID extracts and validates the `:id` URL param into a Mongo ObjectID.
func parseObjectID(c *gin.Context) (primitive.ObjectID, bool) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return primitive.NilObjectID, false
	}
	return objID, true
}

// handleValidationError processes validation errors from JSON binding.
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

// Register handles POST /register
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := data.Register(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User Registered successfully"})
}

// Login handles POST /login.
func Login(ctx *gin.Context) {
	var credentials models.User
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := data.Login(credentials.Username, credentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expire := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// GetTasks handles GET /api/tasks.
func GetTasks(ctx *gin.Context) {
	tasks, err := data.GetTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// GetTaskById handles GET /api/tasks/:id.
func GetTaskById(ctx *gin.Context) {
	objID, ok := parseObjectID(ctx)
	if !ok {
		return
	}

	task, err := data.GetTaskById(objID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// CreateTask handles POST /api/tasks.
func CreateTask(ctx *gin.Context) {
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		handleValidationError(ctx, err)
		return
	}
	created, err := data.CreateTask(newTask)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "task already exists"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

// UpdateTask handles PUT /api/tasks/:id.
func UpdateTask(ctx *gin.Context) {
	objID, ok := parseObjectID(ctx)
	if !ok {
		return
	}

	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		handleValidationError(ctx, err)
		return
	}
	updated, err := data.UpdateTask(objID, updatedTask)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "duplicate value for unique field"})
		} else if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// DeleteTask handles DELETE /api/tasks/:id.
func DeleteTask(ctx *gin.Context) {
	objID, ok := parseObjectID(ctx)
	if !ok {
		return
	}

	err := data.DeleteTask(objID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
