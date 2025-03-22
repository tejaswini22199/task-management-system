package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-service/internal/repository"
	"task-service/internal/utils"
)

// GetTasksHandler - Only fetches tasks for the logged-in user
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromToken(r) // Extract user ID from JWT
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 20
	offset := 0
	queryParams := r.URL.Query()
	if page, exists := queryParams["page"]; exists {
		pageNum, err := strconv.Atoi(page[0])
		if err == nil && pageNum > 0 {
			offset = (pageNum - 1) * limit
		}
	}

	statusFilter := ""
	if status, exists := queryParams["status"]; exists {
		statusFilter = status[0]
	}

	tasks, err := repository.GetTasksForUser(utils.DB, userID, limit, offset, statusFilter)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}
