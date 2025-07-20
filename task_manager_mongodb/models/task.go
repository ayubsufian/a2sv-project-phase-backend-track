package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a single task item in the task manager system.
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	DueDate     time.Time          `json:"due_date" binding:"required,duedate" bson:"due_date"`
	Status      string             `json:"status" binding:"required,oneof=pending completed" bson:"status"`
}
