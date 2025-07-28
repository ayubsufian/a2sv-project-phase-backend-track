package routers

import (
	"task_manager_clean/delivery/controllers"
	"task_manager_clean/infrastructure"

	"github.com/gin-gonic/gin"
)

// SetupRouter constructs the Gin engine with all application routes.
func SetupRouter(h *controllers.Handler, jwtSvc infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// Public routes for registration and login functionality.
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Protected API routes require a valid JWT.
	api := r.Group("/api")
	api.Use(infrastructure.AuthMiddleware(jwtSvc))
	{
		api.GET("/tasks", h.GetTasks)
		api.POST("/tasks", h.CreateTask)
		api.GET("/tasks/:id", h.GetTask)
		api.PUT("/tasks/:id", h.UpdateTask)
		api.DELETE("/tasks/:id", h.DeleteTask)

		// Admin-only subgroup for dashboard access.
		admin := api.Group("/admin")
		admin.Use(infrastructure.AdminOnly())
		admin.GET("/dashboard", h.AdminDashboard)
	}

	return r
}
