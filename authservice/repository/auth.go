package repository

import (
	"database/sql"

	"github.com/lib/pq"
	db "github.com/tejaswini22199/task-management-system/database"
)

func InsertUser(name, email, password string) (int, error) {
	var userID int
	err := db.DB.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		name, email, password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}
	return userID, nil
}

func GetUserByEmail(email string) (int, string, error) {
	var userID int
	var storedPassword string
	err := db.DB.QueryRow(
		"SELECT id, password FROM users WHERE email = $1",
		email,
	).Scan(&userID, &storedPassword)

	if err == sql.ErrNoRows {
		return 0, "", err
	}

	return userID, storedPassword, err
}

func GetExistingUserIDs(userIDs []int) ([]int, error) {
	query := `SELECT id FROM users WHERE id = ANY($1)`

	rows, err := db.DB.Query(query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var existingUserIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existingUserIDs = append(existingUserIDs, id)
	}

	return existingUserIDs, nil
}
