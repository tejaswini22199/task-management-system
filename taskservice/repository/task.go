package repository

import (
	"database/sql"
	"github.com/tejaswini22199/task-management-system/taskservice/models"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func CreateTask(task models.Task) (models.Task, error) {
	query := `INSERT INTO tasks (title, description, status, created_by, created_at) 
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	err := db.QueryRow(query, task.Title, task.Description, task.Status, task.CreatedBy).Scan(&task.ID)
	return task, err
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
