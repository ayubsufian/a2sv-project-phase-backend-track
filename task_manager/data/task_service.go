package data

import (
	"errors"
	"sync"
	"task_manager/models"
	"time"
)

var (
	tasks = []models.Task{
		{
			ID:          1,
			Title:       "Buy groceries",
			Description: "Milk, eggs, bread, and fruits",
			DueDate:     time.Date(2025, 7, 20, 10, 0, 0, 0, time.UTC),
			Status:      "pending",
		},
		{
			ID:          2,
			Title:       "Finish report",
			Description: "Complete monthly sales analysis",
			DueDate:     time.Date(2025, 7, 22, 15, 30, 0, 0, time.UTC),
			Status:      "pending",
		},
		{
			ID:          3,
			Title:       "Team meeting",
			Description: "Discuss project milestones",
			DueDate:     time.Date(2025, 7, 25, 9, 0, 0, 0, time.UTC),
			Status:      "completed",
		},
	}
	// In-memory task list
	nextID = 4        // Auto-incrementing task ID
	mu     sync.Mutex // Mutex for concurrent access
)

// GetAll returns all tasks
func GetAll() []models.Task {
	mu.Lock()
	defer mu.Unlock()
	return tasks
}

// GetByID retrieves a task by ID
func GetByID(id int) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Task{}, errors.New("not found")
}

// Create adds a new task
func Create(t models.Task) models.Task {
	mu.Lock()
	defer mu.Unlock()
	t.ID = nextID
	nextID++
	tasks = append(tasks, t)
	return t
}

// Update modifies an existing task
func Update(id int, t models.Task) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			t.ID = id
			tasks[i] = t
			return t, nil
		}
	}
	return models.Task{}, errors.New("not found")
}

// Delete removes a task by ID
func Delete(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
