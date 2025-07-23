package data

import (
	"context"
	"errors"
	"fmt"
	"task_manager_auth/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// client is the MongoDB client instance used throughout the application.
var client *mongo.Client

// taskCollection holds a reference to the "tasks" collection in the "taskdb" database.
var taskCollection *mongo.Collection

// InitMongoDB initializes the MongoDB connection using the provided URI.
func InitMongoDB(ctx context.Context, uri string) error {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	taskCollection = client.Database("taskdb").Collection("tasks")
	return nil
}

// CloseMongoDB cleanly disconnects the MongoDB client when the application shuts down.
func CloseMongoDB() {
	if client != nil {
		_ = client.Disconnect(context.TODO())
	}
}

// GetTasks retrieves all task documents from the "tasks" collection.
func GetTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTaskById retrieves a single task by its ObjectID.
func GetTaskById(id primitive.ObjectID) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err := taskCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return task, errors.New("task not found")
	}
	return task, err
}

// CreateTask inserts a new task into the "tasks" collection.
func CreateTask(task models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.ID = primitive.NewObjectID()
	_, err := taskCollection.InsertOne(ctx, task)
	return task, err
}

// UpdateTask replaces an existing task document with new data.
func UpdateTask(id primitive.ObjectID, updatedTask models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	updatedTask.ID = id

	res, err := taskCollection.ReplaceOne(ctx, bson.M{"_id": id}, updatedTask)
	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return models.Task{}, fmt.Errorf("duplicate key: %w", err)
			}
		}
	}
	if err != nil {
		return models.Task{}, err
	}
	if res.MatchedCount == 0 {
		return models.Task{}, errors.New("not found")
	}
	return updatedTask, nil
}

// DeleteTask removes a task document by its ObjectID.
func DeleteTask(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := taskCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
