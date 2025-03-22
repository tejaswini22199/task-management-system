package models

import "time"

type TaskUser struct {
    TaskID    int    `json:"task_id"`
    UserID    int    `json:"user_id"`
    // Optionally, you might want to add:
    AssignedAt time.Time `json:"assigned_at"`
}