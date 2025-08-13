package domain

import "time"

// Task represents a user’s to-do item in the system.
type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}
