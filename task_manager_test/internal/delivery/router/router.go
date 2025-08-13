package router

import (
	"task_manager_test/internal/delivery/controller"
	"task_manager_test/internal/delivery/middleware"
	"task_manager_test/internal/usecase"

	"github.com/gin-gonic/gin"
)

// RouterConfig holds the dependencies for the router.
type RouterConfig struct {
	UserCont *controller.UserController
	TaskCont *controller.TaskController
	JwtSvc   usecase.IJWTService
}

// SetupRouter constructs the Gin engine with all application routes.
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	// Public routes for registration and login functionality.
	r.POST("/register", cfg.UserCont.Register)
	r.POST("/login", cfg.UserCont.Login)

	// Protected API routes require a valid JWT.
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JwtSvc))
	{
		api.GET("/tasks", cfg.TaskCont.GetTasks)
		api.POST("/tasks", cfg.TaskCont.CreateTask)
		api.GET("/tasks/:id", cfg.TaskCont.GetTask)
		api.PUT("/tasks/:id", cfg.TaskCont.UpdateTask)
		api.DELETE("/tasks/:id", cfg.TaskCont.DeleteTask)

		// Admin-only subgroup for dashboard access.
		admin := api.Group("/admin")
		admin.Use(middleware.AdminOnly())
		admin.GET("/dashboard", cfg.TaskCont.AdminDashboard)
	}

	return r
}
