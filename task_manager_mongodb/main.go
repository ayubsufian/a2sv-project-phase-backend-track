package main

import (
	"context"
	"log"
	"os"
	"task_manager_mongodb/data"
	"task_manager_mongodb/router"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables from .env file, if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, expecting env vars")
	}

	// Get MongoDB connection URI from environment variables
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	// Create a context with timeout for MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Declare MongoDB client and connection error
	var client *mongo.Client
	var err error
	const maxConnectAttempts = 5 // Maximum retry attempts
	baseDelay := time.Second     // Base delay between retries

	// Attempt to connect to MongoDB with exponential backoff
	for i := 0; i < maxConnectAttempts; i++ {
		client, err = mongo.Connect(ctx, options.Client().
			ApplyURI(uri).
			SetServerSelectionTimeout(15*time.Second).
			SetRetryWrites(true))
		if err == nil {
			err = client.Ping(ctx, nil) // Check if the connection is alive
		}
		if err == nil {
			break // Exit loop if connection successful
		}
		// Log the failed attempt and wait before retrying
		delay := baseDelay * (1 << i)
		log.Printf("MongoDB connection attempt %d failed: %v — retrying in %s", i+1, err, delay)
		time.Sleep(delay)
	}
	// If all attempts fail, terminate the program
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB after retries: %v", err)
	}

	// Assign the connected client and collection to the global data package
	data.Client = client
	data.TasksCollection = client.Database("taskdb").Collection("tasks")

	// Register custom validation for 'duedate' to ensure it's a future date
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("duedate", func(fl validator.FieldLevel) bool {
			date, ok := fl.Field().Interface().(time.Time)
			return ok && date.After(time.Now())
		})
	}

	// Set up and run the Gin router on port 8080
	r := router.SetupRouter()
	r.Run(":8080")
}
