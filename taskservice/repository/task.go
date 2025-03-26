package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/tejaswini22199/task-management-system/database" // Import the db package
	"github.com/tejaswini22199/task-management-system/taskservice/models"
)

// var db *sql.DB

// Initialize the repository with the singleton DB instance
// func InitRepository() {
// 	if db != nil {
// 		log.Println("Repository already initialized")
// 		return
// 	}
// 	db = database.GetDB()
// 	if db == nil {
// 		log.Fatal("Failed to initialize repository: database instance is nil")
// 	} else {
// 		log.Println("✅ Repository initialized successfully")
// 	}
// }

func CreateTask(task models.Task, userIDs []int) (models.Task, error) {
	fmt.Println("line 24")
	db := database.GetDB()
	tx, err := db.Begin() // Start transaction
	fmt.Println("line 32")
	if err != nil {
		fmt.Println("line 34")
		return task, err
	}

	// Insert into tasks table
	query := `INSERT INTO tasks (title, description, status, created_by, created_at) 
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	err = tx.QueryRow(query, task.Title, task.Description, task.Status, task.CreatedBy).Scan(&task.ID)
	if err != nil {
		tx.Rollback() // Rollback on failure
		return task, err
	}

	// Insert into tasks_users table
	if len(userIDs) > 0 {
		err = AssignUsersToTask(tx, task.ID, userIDs)
		if err != nil {
			tx.Rollback() // Rollback if assigning users fails
			return task, err
		}
	}

	err = tx.Commit() // Commit if everything is successful
	return task, err
}

func AssignUsersToTask(tx *sql.Tx, taskID int, userIDs []int) error {
	query := `INSERT INTO tasks_users (task_id, user_id) VALUES `

	placeholders := []string{}
	values := []interface{}{taskID}

	for i, userID := range userIDs {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, userID)
	}

	query += strings.Join(placeholders, ",") // Construct final query

	_, err := tx.Exec(query, values...) // Execute query
	return err
}

func GetTaskIds(userId int) ([]int, error) {
	db := database.GetDB()
	query := "SELECT task_id FROM tasks_users WHERE user_id = $1"
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskIds []int
	for rows.Next() {
		var taskId int
		if err := rows.Scan(&taskId); err != nil {
			return nil, err
		}
		taskIds = append(taskIds, taskId)
	}
	fmt.Println(taskIds)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return taskIds, nil
}

func GetTaskByID(id int) (models.Task, error) {
	fmt.Println("Fetching task with ID:", id)
	var task models.Task

	db := database.GetDB()
	err := db.QueryRow("SELECT id, title, description, status, created_by, created_at FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedBy, &task.CreatedAt)

	if err != nil {
		fmt.Println("Error fetching task:", err)
	}
	return task, err
}

func GetTaskIdsWithPagination(userId, limit, offset int) ([]int, error) {
	fmt.Println("Fetching paginated task IDs for user:", userId)

	db := database.GetDB()
	query := `
		SELECT task_id FROM tasks_users
		WHERE user_id = $1
		ORDER BY task_id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskIds []int
	for rows.Next() {
		var taskId int
		if err := rows.Scan(&taskId); err != nil {
			continue
		}
		taskIds = append(taskIds, taskId)
	}

	return taskIds, nil
}

func GetTasksForUser(userId, page, limit int) ([]models.Task, int, error) {
	fmt.Println("Fetching tasks for user:", userId)

	db := database.GetDB()

	// Calculate OFFSET
	offset := (page - 1) * limit

	// Get total count of tasks for user (for pagination)
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks_users WHERE user_id = $1", userId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated tasks
	query := `
		SELECT t.id, t.title, t.description, t.status, t.created_at
		FROM tasks t
		JOIN tasks_users tu ON t.id = tu.task_id
		WHERE tu.user_id = $1
		ORDER BY id ASC 
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			continue // Skip if a task can't be scanned
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}
func GetTasksByStatus(status string, userId int) ([]models.Task, error) {

	// Step 1: Fetch task IDs for the given user
	taskIds, err := GetTaskIds(userId)
	fmt.Println(taskIds)
	if err != nil {
		fmt.Println("Line 192")
		return nil, err
	}

	// Step 2: If no tasks are assigned to the user, return an empty list
	if len(taskIds) == 0 {
		fmt.Println("Line 198")
		return []models.Task{}, nil
	}

	// Step 3: Fetch tasks with given status for the obtained task IDs
	tasks, err := FetchTasksByStatusAndIds(status, taskIds)
	if err != nil {
		fmt.Print("line 204")
		return nil, err
	}
	fmt.Println(tasks)
	fmt.Println("line 208")
	return tasks, nil
}

// // FetchTasksByStatusAndIds retrieves tasks based on status and task IDs
// func FetchTasksByStatusAndIds(status string, taskIds []int) ([]models.Task, error) {
// 	db := database.GetDB()
// 	// Generate placeholders dynamically: ?, ?, ?
// 	placeholders := strings.Repeat("?, ", len(taskIds))
// 	placeholders = strings.TrimSuffix(placeholders, ", ") // Remove last comma

// 	query := fmt.Sprintf(
// 		"SELECT id, title, description, status, created_at FROM tasks WHERE status = ? AND id IN (%s)",
// 		placeholders,
// 	)

// 	// Prepare arguments correctly: first status, then task IDs
// 	args := append([]interface{}{status}, ConvertToInterfaceSlice(taskIds)...)

// 	rows, err := db.Query(query, args...)
// 	if err != nil {
// 		log.Println("Query Error:", err)
// 	}
// 	defer rows.Close()

// 	fmt.Println("Query Executed Successfully")
// 	rows, err := db.Query("SELECT id, title, description, status, created_at FROM tasks WHERE user_id = $1", userId)
// 	if err != nil {
// 		return nil, err // Return early if query fails
// 	}
// 	defer rows.Close() // Always close the rows after reading
// 	tasks, err := ScanTasks(rows)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ScanTasks(tasks)
// }

// FetchTasksByStatusAndIds retrieves tasks based on status and task IDs
func FetchTasksByStatusAndIds(status string, taskIds []int) ([]models.Task, error) {
	db := database.GetDB()

	// Generate PostgreSQL placeholders: $2, $3, $4
	placeholders := make([]string, len(taskIds))
	args := []interface{}{status} // First argument is status

	for i, id := range taskIds {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // Start from $2
		args = append(args, id)                   // Append IDs to args
	}

	query := fmt.Sprintf(
		"SELECT id, title, description, status, created_at FROM tasks WHERE status = $1 AND id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GeneratePlaceholders returns a string of `?` placeholders for SQL queries
func GeneratePlaceholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?,", count), ",")
}

// ConvertToInterfaceSlice converts an int slice to an interface{} slice for SQL queries
func ConvertToInterfaceSlice(ints []int) []interface{} {
	intfSlice := make([]interface{}, len(ints))
	for i, v := range ints {
		intfSlice[i] = v
	}
	return intfSlice
}

func ScanTasks(rows *sql.Rows) ([]models.Task, error) {
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, err // Return the error if scanning fails
		}
		tasks = append(tasks, task)
	}

	// Check if there was an error during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println("line 260")
	fmt.Println(tasks)
	return tasks, nil
}

func UpdateTask(task models.Task) (models.Task, error) {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4 RETURNING id, title, description, status, created_at`
	db := database.GetDB()

	// Scan the updated task to return it
	var updatedTask models.Task
	err := db.QueryRow(query, task.Title, task.Description, task.Status, task.ID).
		Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.CreatedAt)

	if err != nil {
		return models.Task{}, err
	}

	return updatedTask, nil
}

func DeleteTask(id int) error {
	db := database.GetDB()

	// Execute the DELETE query and get the number of affected rows
	result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Check if any row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
