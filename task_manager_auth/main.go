package main

import (
	"context"
	"log"
	"os"
	"task_manager_auth/data"
	"task_manager_auth/router"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from a .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create a root context with a 5-second timeout to use during DB initialization
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the MongoDB connection URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	// Connect to MongoDB using the context with timeout to prevent hanging
	err := data.InitMongoDB(ctx, mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	// Ensure the connection is properly closed when the application exits
	defer data.CloseMongoDB()

	// Set up the HTTP router and start the server
	data.InitUserCollection()
	r := router.SetUpRouter()
	r.Run(":8080")
}
