package models

import "time"

// TaskStatus enum
type TaskStatus string

const (
	ToDo       TaskStatus = "To Do"       // Task is planned but not started
	InProgress TaskStatus = "In Progress" // Work on the task has begun
	Completed  TaskStatus = "Completed"   // Task is finished
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      TaskStatus   `json:"status"`
	CreatedBy   int        `json:"created_by"`  // User ID of the task creator 
	CreatedAt   time.Time `json:"created_at"` // List of user IDs assigned to this task
}

