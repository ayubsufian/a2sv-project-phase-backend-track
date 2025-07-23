package router

import (
	"task_manager_auth/controllers"
	"task_manager_auth/middleware"

	"github.com/gin-gonic/gin"
)

// SetUpRouter configures all HTTP routes, groups, and middlewares.
func SetUpRouter() *gin.Engine {
	// Create a Gin router with default middleware.
	r := gin.Default()

	// Public routes for user registration and login
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Create a /api group which will require authenticated access
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/tasks", controllers.GetTasks)
		api.POST("/tasks", controllers.CreateTask)
		api.GET("/tasks/:id", controllers.GetTaskById)
		api.PUT("/tasks/:id", controllers.UpdateTask)
		api.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	// Nested /api/admin group requiring admin privileges
	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome Admin"})
		})
	}

	// Return the configured router ready to be run
	return r
}
