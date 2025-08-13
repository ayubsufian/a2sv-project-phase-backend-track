package main

import (
	"context"
	"log"
	"os"
	"task_manager_test/internal/delivery/controller"
	"task_manager_test/internal/delivery/router"
	"task_manager_test/internal/repository"
	"task_manager_test/internal/service"
	"task_manager_test/internal/usecase"
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

	// Get a handle to the database.
	db := client.Database("taskdb")

	// Initialize repositories with the database handle.
	taskRepo := repository.NewMongoTaskRepository(db)
	userRepo := repository.NewMongoUserRepository(db)

	// Initialize services with the correct types.
	pwdSvc := service.NewPasswordHasher()
	jwtSvc := service.NewJWTService(jwtSecret)

	// Initialize usecases (business logic) for users and tasks.
	userUC := usecase.NewUserUsecase(userRepo, pwdSvc, jwtSvc)
	taskUC := usecase.NewTaskUsecase(taskRepo)

	// Initialize each controller individually.
	userCont := controller.NewUserController(userUC)
	taskCont := controller.NewTaskController(taskUC)

	// Populate the RouterConfig struct
	routerCfg := &router.RouterConfig{
		UserCont: userCont,
		TaskCont: taskCont,
		JwtSvc:   jwtSvc,
	}

	// Set up the HTTP router with the config struct.
	router := router.SetupRouter(routerCfg)

	// Start the HTTP server on port 8080.
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
