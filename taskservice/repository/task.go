package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tejaswini22199/task-management-system/database" // Import the db package
	"github.com/tejaswini22199/task-management-system/taskservice/models"
)

var db *sql.DB

// Initialize the repository with the singleton DB instance
func InitRepository() {
	db = database.GetDB() // Fetch the DB instance from the singleton
}

func CreateTask(task models.Task, userIDs []int) (models.Task, error) {
	tx, err := db.Begin() // Start transaction

	if err != nil {
		fmt.Println("line 23")
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

func GetTasks() ([]models.Task, error) {
	rows, err := db.Query("SELECT id, title, description, status, created_at FROM tasks")
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

func GetTaskByID(id int) (models.Task, error) {
	var task models.Task
	err := db.QueryRow("SELECT id, title, description, status, created_at FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	return task, err
}

func GetTasksByStatus(status string) ([]models.Task, error) {
	rows, err := db.Query("SELECT id, title, description, status, created_at FROM tasks WHERE status = $1", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func UpdateTask(task models.Task) (models.Task, error) {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4`
	_, err := db.Exec(query, task.Title, task.Description, task.Status, task.ID)
	return task, err
}

func DeleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
