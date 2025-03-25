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

	fmt.Println("line 59")

	validUserIDs, isValid := utils.ValidateInputUserIDs(input.UserIDs, c)

	if !isValid {
		return // Stop execution if user ID validation fails
	}

	fmt.Println("line 67")
	task, err := services.CreateTask(input, userID, validUserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task, "message": "Task created successfully"})
}

func GetTasks(c *gin.Context) {
	tasks, err := services.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTaskByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	task, err := services.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func GetTasksByStatus(c *gin.Context) {
	status := c.Param("status")

	tasks, err := services.GetTasksByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func UpdateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	task, err := services.UpdateTask(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task, "message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := services.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
