package repository

import (
	db "github.com/tejaswini22199/task-management-system/database"
	"database/sql"
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
