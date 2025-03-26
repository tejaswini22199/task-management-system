package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tejaswini22199/task-management-system/authservice/utils"
	"github.com/tejaswini22199/task-management-system/taskservice/models"
	"github.com/tejaswini22199/task-management-system/taskservice/services"
	taskUtils "github.com/tejaswini22199/task-management-system/taskservice/utils"
)

// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	// Convert userID to int (if needed)
// 	authUserID, ok := userID.(int)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

func CreateTask(c *gin.Context) {

	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Authenticated User ID:", userID)

	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("the value of input status " + input.Status)

	isValid, errMsg := taskUtils.ValidateTaskStatus(input.Status)

	if !isValid {
		fmt.Println("task status is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	validUserIDs, isValid := utils.ValidateInputUserIDs(input.UserIDs, c)

	if !isValid {
		return // Stop execution if user ID validation fails
	}

	task, err := services.CreateTask(input, userID, validUserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task, "message": "Task created successfully"})
}

func GetTasks(c *gin.Context) {
	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Extract pagination parameters from query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default page = 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20")) // Default limit = 20

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	fmt.Println("Authenticated User ID:", userID)
	tasks, total, err := services.GetTasks(userID, page, limit) // Pass page & limit
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total_tasks": total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + limit - 1) / limit, // Calculate total pages
	})
}

func GetTaskByID(c *gin.Context) {
	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := services.GetTaskByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func GetTasksByStatus(c *gin.Context) {
	status := c.Param("status")
	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tasks, err := services.GetTasksByStatus(status, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func UpdateTask(c *gin.Context) {
	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"user_id": userID, // Optional: Might be nil in case of an error
		})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"user_id": userID,
		})
		return
	}

	task, err := services.UpdateTask(input.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task, "message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	userID, err := utils.ValidateUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = services.DeleteTask(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
