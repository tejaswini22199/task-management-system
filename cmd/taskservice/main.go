// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"github.com/gin-gonic/gin"
// 	"github.com/dgrijalva/jwt-go"
// 	"database/sql"
// 	"github.com/lib/pq"
// )

// var db *sql.DB
// const AuthServiceURL = "http://localhost:8000/validate"

// func main() {
// 	log.Println("Starting Auth Service...")
// 	db.InitDB()

// 	r := gin.Default()
// 	r.Use(authMiddleware) // Protects all task routes

// 	r.POST("/tasks", createTask)
// 	r.GET("/tasks", getTasks)
// 	r.GET("/tasks/status/:status", getTasksByStatus)
// 	r.GET("/tasks/:id", getTaskByID)
// 	r.PUT("/tasks/:id", updateTask)
// 	r.DELETE("/tasks/:id", deleteTask)

// 	log.Println("Task Service running on port 8001")
// 	r.Run(":8001")
// }

// // Auth Middleware
// func authMiddleware(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
// 		c.Abort()
// 		return
// 	}

// 	bearerToken := strings.Split(authHeader, " ")
// 	if len(bearerToken) != 2 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 		c.Abort()
// 		return
// 	}

// 	tokenString := bearerToken[1]
// 	claims := jwt.MapClaims{}

// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("your-secret-key"), nil
// 	})

// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		c.Abort()
// 		return
// 	}

// 	userID := claims["user_id"].(float64)
// 	c.Set("user_id", int(userID))
// 	c.Next()
// }

// func createTask(c *gin.Context) {
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

// 	// Now you can use authUserID to track who created the task
// 	var input struct {
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		Status      string `json:"status"`
// 		UserIDs     []int  `json:"user_ids"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	status := models.TaskStatus(input.Status)

// 	validStatuses := map[models.TaskStatus]bool{
// 		models.ToDo:       true,
// 		models.InProgress: true,
// 		models.Completed:  true,
// 	}

// 	if _, isValid := validStatuses[status]; !isValid {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed values: To Do, In Progress, Completed"})
// 		return
// 	}

// 	// Ensure input user IDs are unique using a set
// 	uniqueUserIDs := make(map[int]struct{})
// 	for _, id := range input.UserIDs {
// 		uniqueUserIDs[id] = struct{}{}
// 	}

// 	// Convert set to slice for database validation
// 	validUserIDs := make([]int, 0, len(uniqueUserIDs))
// 	for id := range uniqueUserIDs {
// 		validUserIDs = append(validUserIDs, id)
// 	}

// 	// Validate if the given user IDs exist in the database
// 	if len(validUserIDs) > 0 {
// 		query := `SELECT id FROM users WHERE id = ANY($1)`
// 		rows, err := db.Query(query, pq.Array(validUserIDs))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		}
// 		defer rows.Close()

// 		// Collect existing users in a set
// 		existingUsers := make(map[int]struct{})
// 		for rows.Next() {
// 			var id int
// 			rows.Scan(&id)
// 			existingUsers[id] = struct{}{}
// 		}

// 		// Identify non-existent users
// 		invalidUserIDs := []int{}
// 		for _, id := range validUserIDs {
// 			if _, exists := existingUsers[id]; !exists {
// 				invalidUserIDs = append(invalidUserIDs, id)
// 			}
// 		}

// 		// If some users don't exist, return an error
// 		if len(invalidUserIDs) > 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error":         "Some user IDs do not exist",
// 				"invalid_users": invalidUserIDs,
// 			})
// 			return
// 		}
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}

// 	// Insert task and get task ID
// 	var taskID int
// 	query := `INSERT INTO tasks (title, description, status, created_by, created_at)
// 			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
// 	err = tx.QueryRow(query, input.Title, input.Description, input.Status, authUserID).Scan(&taskID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
// 		return
// 	}

// 	// Assign task to valid users
// 	if len(validUserIDs) > 0 {
// 		values := []interface{}{taskID}
// 		placeholders := []string{}
// 		for i, uid := range validUserIDs {
// 			placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
// 			values = append(values, uid)
// 		}

// 		query := fmt.Sprintf("INSERT INTO tasks_users (task_id, user_id) VALUES %s", strings.Join(placeholders, ","))
// 		_, err = tx.Exec(query, values...)
// 		if err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to users"})
// 			return
// 		}
// 	}

// 	tx.Commit()

// 	c.JSON(http.StatusCreated, gin.H{"task_id": taskID, "assigned_users": validUserIDs, "message": "Task created successfully"})
// }

// func getTaskByID(c *gin.Context) {
// 	taskID := c.Param("id")
// 	var task models.Task

// 	query := "SELECT id, title, description, status, created_at FROM tasks WHERE id = $1"
// 	err := db.QueryRow(query, taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
// 		return
// 	}

// 	var userIDs []int
// 	rows, err := db.Query("SELECT user_id FROM tasks_users WHERE task_id = $1", taskID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned users"})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var userID int
// 		if err := rows.Scan(&userID); err == nil {
// 			userIDs = append(userIDs, userID)
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"task": gin.H{
// 		"id":          task.ID,
// 		"title":       task.Title,
// 		"description": task.Description,
// 		"status":      task.Status,
// 		"created_at":  task.CreatedAt,
// 		"user_ids":    userIDs,
// 	}})
// }

// func getTasks(c *gin.Context) {
// 	var tasks []models.Task
// 	var total int

// 	// Default pagination values
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

// 	// Ensure valid values
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}

// 	offset := (page - 1) * limit

// 	// Get total count of tasks
// 	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total task count"})
// 		return
// 	}

// 	// Fetch paginated tasks
// 	query := "SELECT id, title, description, status, created_at FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2"
// 	rows, err := db.Query(query, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
// 		return
// 	}
// 	defer rows.Close()

// 	// Map for storing taskID to userIDs
// 	taskUsers := make(map[int][]int)

// 	// Retrieve tasks
// 	for rows.Next() {
// 		var task models.Task
// 		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tasks"})
// 			return
// 		}
// 		tasks = append(tasks, task)
// 	}

// 	// Fetch user assignments
// 	userRows, err := db.Query("SELECT task_id, user_id FROM tasks_users WHERE task_id IN (SELECT id FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2)", limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned users"})
// 		return
// 	}
// 	defer userRows.Close()

// 	for userRows.Next() {
// 		var taskID, userID int
// 		if err := userRows.Scan(&taskID, &userID); err == nil {
// 			taskUsers[taskID] = append(taskUsers[taskID], userID)
// 		}
// 	}

// 	// Attach assigned users to tasks
// 	var taskList []gin.H
// 	for _, task := range tasks {
// 		taskList = append(taskList, gin.H{
// 			"id":          task.ID,
// 			"title":       task.Title,
// 			"description": task.Description,
// 			"status":      task.Status,
// 			"created_at":  task.CreatedAt,
// 			"user_ids":    taskUsers[task.ID], // Get user IDs for the task
// 		})
// 	}

// 	// Pagination metadata
// 	c.JSON(http.StatusOK, gin.H{
// 		"tasks":       taskList,
// 		"total_tasks": total,
// 		"page":        page,
// 		"limit":       limit,
// 		"total_pages": (total + limit - 1) / limit, // Calculate total pages
// 	})
// }

// func getTasksByStatus(c *gin.Context) {
// 	status := c.Param("status")

// 	rows, err := db.Query("SELECT id, title, description, status, created_at FROM tasks WHERE status = $1", status)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
// 		return
// 	}
// 	defer rows.Close()

// 	var tasks []map[string]interface{}
// 	for rows.Next() {
// 		var task map[string]interface{}
// 		var id int
// 		var title, description, taskStatus string
// 		var createdAt time.Time

// 		err := rows.Scan(&id, &title, &description, &taskStatus, &createdAt)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning task data"})
// 			return
// 		}

// 		task = map[string]interface{}{
// 			"id":          id,
// 			"title":       title,
// 			"description": description,
// 			"status":      taskStatus,
// 			"created_at":  createdAt,
// 		}

// 		tasks = append(tasks, task)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
// }

// func checkTaskOwnership(c *gin.Context, taskID string) (bool, int, error) {
// 	// Get the task's owner from the database
// 	var ownerID int
// 	err := db.QueryRow("SELECT user_id FROM tasks_users WHERE task_id = $1", taskID).Scan(&ownerID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, 0, fmt.Errorf("task not found")
// 		}
// 		return false, 0, fmt.Errorf("failed to fetch task")
// 	}

// 	// Retrieve the user ID from the context
// 	userID := c.MustGet("user_id").(int)

// 	// Check if the current user is the task owner
// 	if userID != ownerID {
// 		return false, 0, fmt.Errorf("you are not authorized to access or modify this task")
// 	}

// 	// Return true if the user is the owner
// 	return true, ownerID, nil
// }

// func updateTask(c *gin.Context) {
// 	id := c.Param("id")
// 	var task models.Task

// 	// Check if the user is authorized to update this task
// 	authorized, _, err := checkTaskOwnership(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if !authorized {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this task"})
// 		return
// 	}

// 	// Bind the incoming JSON data to the task model
// 	if err := c.ShouldBindJSON(&task); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Update query to modify the task
// 	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4 RETURNING id, title, description, status, created_at"
// 	err = db.QueryRow(query, task.Title, task.Description, task.Status, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)

// 	// Check if the task was updated successfully
// 	if err != nil {
// 		// If no rows were updated, the task ID does not exist in the database
// 		if err.Error() == "no rows in result set" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
// 			return
// 		}
// 		// For any other errors, send a generic failure message
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
// 		return
// 	}

// 	// Return the updated task
// 	c.JSON(http.StatusOK, task)
// }

// func deleteTask(c *gin.Context) {
// 	id := c.Param("id")

// 	// Check if the user is authorized to delete this task
// 	authorized, _, err := checkTaskOwnership(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if !authorized {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this task"})
// 		return
// 	}

// 	// Check if the task exists before proceeding with deletion
// 	var taskID int
// 	err = db.QueryRow("SELECT id FROM tasks WHERE id = $1", id).Scan(&taskID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check task existence"})
// 		}
// 		return
// 	}

// 	// If the task exists, proceed with deletion
// 	tx, err := db.Begin()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
// 		return
// 	}

// 	// Remove task-user mapping
// 	_, err = tx.Exec("DELETE FROM tasks_users WHERE task_id = $1", id)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task-user mappings"})
// 		return
// 	}

// 	// Delete the task itself
// 	res, err := tx.Exec("DELETE FROM tasks WHERE id = $1", id)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
// 		return
// 	}
// 	tx.Commit()

// 	rowsAffected, _ := res.RowsAffected()
// 	if rowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK 1234
// 		, gin.H{"message": "Task deleted successfully"})
// }

//  func authMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}

// 		bearerToken := strings.Split(authHeader, " ")
// 		if len(bearerToken) != 2 {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := bearerToken[1]
// 		claims := &Claims{}

// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(JWTSecret), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("user_id", claims.UserID)
// 		c.Next()
// 	}
// }

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	database "github.com/tejaswini22199/task-management-system/database"
	"github.com/tejaswini22199/task-management-system/taskservice/repository"
	"github.com/tejaswini22199/task-management-system/taskservice/routes"
)

func main() {
	log.Println("Starting Tasks Service...")

	// Initialize Database Connection
	database.InitDB()
	repository.InitRepository()

	// Create a new Gin router
	r := gin.Default()

	// Register routes with dependencies
	routes.RegisterTaskRoutes(r)

	log.Println("Tasks Service is running on port 8001")
	r.Run(":8001")
}
