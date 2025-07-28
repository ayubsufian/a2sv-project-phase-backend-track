package main

import (
	"context"
	"log"
	"os"
	"task_manager_clean/delivery/controllers"
	"task_manager_clean/delivery/routers"
	"task_manager_clean/infrastructure"
	"task_manager_clean/repositories"
	"task_manager_clean/usecases"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables from a .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read MongoDB connection URI from environment variables.
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	// Read JWT secret key for authentication from environment variables.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	// Create a context with a timeout for MongoDB connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB using the provided URI.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the MongoDB client disconnects when the application stops.
	defer client.Disconnect(ctx)

	// Get references to the "tasks" and "users" collections in the "taskdb" database.
	taskCol := client.Database("taskdb").Collection("tasks")
	userCol := client.Database("taskdb").Collection("users")

	// Initialize repositories for tasks and users.
	taskRepo := repositories.NewMongoTaskRepository(taskCol)
	userRepo := repositories.NewMongoUserRepository(userCol)

	// Initialize infrastructure services: password hasher and JWT service.
	pwdSvc := infrastructure.NewPasswordHasher()
	jwtSvc := infrastructure.NewJWTService([]byte(jwtSecret))

	// Initialize usecases (business logic) for users and tasks.
	userUC := usecases.NewUserUsecase(userRepo, pwdSvc, jwtSvc)
	taskUC := usecases.NewTaskUsecase(taskRepo)

	// Initialize HTTP handlers (controllers) with the usecases.
	handler := controllers.NewHandler(userUC, taskUC)

	// Set up the HTTP router and apply middleware (like JWT authentication).
	router := routers.SetupRouter(handler, jwtSvc)

	// Start the HTTP server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
