package router

import (
	"task_manager_mongodb/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes all route handlers for the API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	tasks := r.Group("/tasks")
	{
		tasks.GET("", controllers.GetTasks)          // GET /tasks
		tasks.GET("/:id", controllers.GetTask)       // GET /tasks/:id
		tasks.POST("", controllers.CreateTask)       // POST /tasks
		tasks.PUT("/:id", controllers.UpdateTask)    // PUT /tasks/:id
		tasks.DELETE("/:id", controllers.DeleteTask) // DELETE /tasks/:id
	}

	return r
}
