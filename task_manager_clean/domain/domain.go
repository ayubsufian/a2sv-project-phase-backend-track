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

// User represents an account in the system.
type User struct {
	ID       string
	Username string
	Password string
	Role     string
}
