package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tejaswini22199/task-management-system/internal/models"
)

var db *sql.DB


func init() {
	var err error

	// Fetch values from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST") // Use this if DB is in another container
	port := os.Getenv("DB_PORT")

	if host == "" {
		host = "localhost" // Default to localhost if not provided
	}

	// PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%ssslmode=disable", host, user, password, dbname, port)

	// Open database connection
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
}

func createTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Status      string `json:"status" binding:"required"`
		UserIDs     []int  `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Println("[ERROR] JSON binding error:", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("[ERROR] Failed to start transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var taskID int
	query := `INSERT INTO tasks (title, description, status, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
	err = tx.QueryRow(query, input.Title, input.Description, input.Status).Scan(&taskID)
	if err != nil {
		tx.Rollback()
		log.Println("[ERROR] Task insertion failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	for _, userID := range input.UserIDs {
		_, err := tx.Exec("INSERT INTO task_users (task_id, user_id) VALUES ($1, $2)", taskID, userID)
		if err != nil {
			tx.Rollback()
			log.Println("[ERROR] Task-user mapping failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to users"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"task_id": taskID, "message": "Task created successfully"})
}

func getTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	var task models.Task

	query := "SELECT id, title, description, status, created_at FROM tasks WHERE id = $1"
	err := db.QueryRow(query, taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var userIDs []int
	rows, err := db.Query("SELECT user_id FROM task_users WHERE task_id = $1", taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err == nil {
			userIDs = append(userIDs, userID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"task": gin.H{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"created_at":  task.CreatedAt,
		"user_ids":    userIDs,
	}})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4 RETURNING id, title, description, status, created_at"
	err := db.QueryRow(query, task.Title, task.Description, task.Status, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	_, err = tx.Exec("DELETE FROM task_users WHERE task_id = $1", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task-user mappings"})
		return
	}

	res, err := tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	tx.Commit()

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func main() {
	r := gin.Default()
	r.POST("/tasks", createTask)
	r.GET("/tasks/:id", getTaskByID)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	log.Println("Task Management Service running on port 8080")
	r.Run(":8080")
}
