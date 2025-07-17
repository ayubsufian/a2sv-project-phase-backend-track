package main

import (
	"task_manager/router"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Register custom validation: due date must be in the future
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("duedate", func(fl validator.FieldLevel) bool {
			date, ok := fl.Field().Interface().(time.Time)
			return ok && date.After(time.Now())
		})
	}

	// Initialize and run the router
	r := router.SetupRouter()
	r.Run(":8080") // Listen and serve on port 8080
}
