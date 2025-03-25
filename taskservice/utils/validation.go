package utils

import (
	"errors"

	"github.com/tejaswini22199/task-management-system/taskservice/models"
)

// ValidateTaskStatus checks if the provided status is valid
func ValidateTaskStatus(status models.TaskStatus) error {
	taskStatus := models.TaskStatus(status)

	validStatuses := map[models.TaskStatus]bool{
		models.ToDo:       true,
		models.InProgress: true,
		models.Completed:  true,
	}

	if _, isValid := validStatuses[taskStatus]; !isValid {
		return errors.New("invalid status. Allowed values: To Do, In Progress, Completed")
	}

	return nil
}
