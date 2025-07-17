package models

import "time"

// Task represents a to-do task with validation tags
type Task struct {
	ID          int       `json:"id"`                                                                          // Auto-assigned ID
	Title       string    `json:"title" binding:"required"`                                                    // Required title
	Description string    `json:"description"`                                                                 // Optional description
	DueDate     time.Time `json:"due_date" binding:"required,duedate" time_format:"2006-01-02T15:04:05Z07:00"` // Required and must be future date
	Status      string    `json:"status" binding:"required,oneof=pending completed"`                           // Required and must be 'pending' or 'completed'
}
