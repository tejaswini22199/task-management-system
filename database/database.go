package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB initializes the database connection only once
func InitDB() {
	once.Do(func() {
		var err error

		// Fetch database credentials from environment variables
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")

		log.Println("Initializing Database...")

		log.Printf("DB Config: host=%s user=%s dbname=%s port=%s\n", host, user, dbname, port)

		if user == "" || password == "" || dbname == "" || host == "" || port == "" {
			log.Fatal("Missing database environment variables")
		}

		// Default to localhost if no host is specified
		if host == "" {
			host = "localhost"
		}

		// Construct the connection string
		connStr := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port,
		)

		// Open a database connection
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Verify the connection
		if err = db.Ping(); err != nil {
			log.Fatalf("Database is unreachable: %v", err)
		}

		log.Println("Connected to PostgreSQL successfully!")

		// Run database migrations
		migrateDB()
	})
}

// GetDB returns the singleton database instance
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return db
}

// migrateDB ensures required tables exist in the database
func migrateDB() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) DEFAULT 'todo',
			created_by INTEGER NOT NULL,  
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS tasks_users (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
			assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Fatalf("Error creating table: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Database tables ensured successfully!")
}
