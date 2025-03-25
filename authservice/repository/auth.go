package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/tejaswini22199/task-management-system/database"
)

var db *sql.DB

// InitRepository initializes the repository with the singleton DB instance
func InitRepository() {
	db = database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize repository: database instance is nil")
	}
}

// InsertUser inserts a new user into the database
func InsertUser(name, email, password string) (int, error) {
	if db == nil {
		return 0, errors.New("database connection is not initialized")
	}

	var userID int
	err := db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		name, email, password,
	).Scan(&userID)

	if err != nil {
		// Handle unique constraint error for duplicate emails
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, errors.New("email already exists")
		}
		return 0, err
	}
	return userID, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (int, string, error) {
	if db == nil {
		return 0, "", errors.New("database connection is not initialized")
	}

	var userID int
	var storedPassword string
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE email = $1",
		email,
	).Scan(&userID, &storedPassword)

	if err == sql.ErrNoRows {
		return 0, "", nil // Return nil error for "not found" cases
	} else if err != nil {
		return 0, "", err
	}

	return userID, storedPassword, nil
}

// GetExistingUserIDs fetches user IDs that exist in the database
func GetExistingUserIDs(userIDs []int) ([]int, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	fmt.Println("line 73")
	query := `SELECT id FROM users WHERE id = ANY($1)`
	rows, err := db.Query(query, pq.Array(userIDs))

	fmt.Println("line 77")

	if err != nil {
		return nil, err
	}

	fmt.Println("line 83")

	defer rows.Close()

	var existingUserIDs []int
	for rows.Next() {
		var id int
		fmt.Println("line 90")
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		fmt.Println("line 94")
		existingUserIDs = append(existingUserIDs, id)
	}

	fmt.Println("line 98")

	return existingUserIDs, nil
}
