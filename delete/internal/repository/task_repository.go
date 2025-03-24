// package repository

// import (
// 	"database/sql"
// 	"log"
// 	"task-service/internal/models"
// )

// func GetTasksForUser(db *sql.DB, userID int, limit int, offset int, statusFilter string) ([]models.Task, error) {
// 	var tasks []models.Task

// 	query := `SELECT t.id, t.title, t.description, t.status, t.created_at 
// 	          FROM tasks t 
// 	          INNER JOIN user_tasks ut ON t.id = ut.task_id 
// 	          WHERE ut.user_id = $1`
	
// 	// Apply filtering
// 	if statusFilter != "" {
// 		query += " AND t.status = $2"
// 		query += " LIMIT $3 OFFSET $4"
// 		rows, err := db.Query(query, userID, statusFilter, limit, offset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rows.Close()
// 	} else {
// 		query += " LIMIT $2 OFFSET $3"
// 		rows, err := db.Query(query, userID, limit, offset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rows.Close()
// 	}

// 	for rows.Next() {
// 		var task models.Task
// 		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
// 			log.Println("Error scanning task:", err)
// 			continue
// 		}
// 		tasks = append(tasks, task)
// 	}

// 	return tasks, nil
// }
