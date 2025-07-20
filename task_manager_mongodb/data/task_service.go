package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"task_manager_mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global MongoDB client and collection references.
var Client *mongo.Client
var TasksCollection *mongo.Collection

// GetAll retrieves all task documents from the MongoDB collection.
func GetAll() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := TasksCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []models.Task
	if err = cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetByID retrieves a task by its ObjectID from the MongoDB collection.
func GetByID(id primitive.ObjectID) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var t models.Task
	err := TasksCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return t, errors.New("not found")
	}
	return t, nil
}

// Create inserts a new task into the MongoDB collection.
func Create(t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.ID = primitive.NewObjectID()
	_, err := TasksCollection.InsertOne(ctx, t)
	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return models.Task{}, fmt.Errorf("duplicate key: %w", err)
			}
		}
	}
	return t, err
}

// Update replaces an existing task document by ID.
func Update(id primitive.ObjectID, t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t.ID = id
	res, err := TasksCollection.ReplaceOne(ctx, bson.M{"_id": id}, t)
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
	return t, nil
}

// Delete removes a task document from the collection by its ID.
func Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := TasksCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || res.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}
