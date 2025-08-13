package controller

import (
	"errors"
	"net/http"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserController wraps use case interfaces for user operations.
type UserController struct {
	userUC usecase.UserUsecase
}

// NewUserController creates a new Handler given User use cases.
func NewUserController(u usecase.UserUsecase) *UserController {
	return &UserController{userUC: u}
}

// Register handles new user registration requests.
func (uc *UserController) Register(c *gin.Context) {
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
	if err := uc.userUC.Register(c.Request.Context(), user); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "a user with this username already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register user"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User Registered successfully"})
}

// Login handles user authentication.
func (uc *UserController) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := uc.userUC.Login(c.Request.Context(), body.Username, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound), errors.Is(err, usecase.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal server error occurred"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
