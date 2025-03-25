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
	ID          int        `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Status      TaskStatus `json:"status" binding:"required"`
	CreatedBy   int        `json:"created_by"` // User ID of the task creator
	CreatedAt   time.Time  `json:"created_at"` // List of user IDs assigned to this task
}

// TaskInput is used for creating or updating a task
type TaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Status      TaskStatus `json:"status" binding:"required"`
	CreatedBy   int        `json:"created_by" binding:"required"`
}
