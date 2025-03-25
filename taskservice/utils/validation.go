package utils

import (
	"github.com/tejaswini22199/task-management-system/taskservice/models"
)

// ValidateTaskStatus checks if the provided status is valid
func ValidateTaskStatus(status models.TaskStatus) (bool, string) {
	validStatuses := map[models.TaskStatus]bool{
		models.ToDo:       true,
		models.InProgress: true,
		models.Completed:  true,
	}

	if _, isValid := validStatuses[status]; !isValid {
		return false, "Invalid status. Allowed values: To Do, In Progress, Completed"
	}

	return true, ""
}
