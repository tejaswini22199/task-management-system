// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"
// 	"github.com/gin-gonic/gin"
// 	"github.com/lib/pq"
// 	"database/sql"
// 	"github.com/golang-jwt/jwt/v5"
// 	"golang.org/x/crypto/bcrypt"
// )

// var db *sql.DB
// const JWTSecret = "your-secret-key"

// // func main() {
// //     log.Println("Starting Auth Service...")
// // 	db.InitDB()
// // 	r := gin.Default()
// // 	r.POST("/register", registerUser)
// // 	r.POST("/login", loginUser)
// // 	log.Println("Auth Service is running on port 8000")
// // 	r.Run(":8000")
// // }

// func main() {
// 	log.Println("Starting Auth Service...")

// 	// Initialize Database Connection
// 	db.InitDB()

// 	// Create a new Gin router
// 	r := gin.Default()

// 	// Register routes
// 	routes.RegisterAuthRoutes(r)

// 	log.Println("Auth Service is running on port 8000")
// 	r.Run(":8000")
// }

// // Generate JWT token
// func generateToken(userID int) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(JWTSecret))
// }

// // Hash password
// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// // Check password hash
// func checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// // Modified register handler with password hashing
// func registerUser(c *gin.Context) {
// 	var input RegisterInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Hash password
// 	hashedPassword, err := hashPassword(input.Password)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
// 		return
// 	}

// 	var userID int
// 	err = db.QueryRow(
// 		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
// 		input.Name, input.Email, hashedPassword,
// 	).Scan(&userID)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	// Generate token for new user
// 	token, err := generateToken(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"user_id": userID,
// 		"token":   token,
// 		"message": "User registered successfully",
// 	})
// }

// // Modified login handler with password verification and JWT
// func loginUser(c *gin.Context) {
// 	var input LoginInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	var userID int
// 	var storedPassword string
// 	err := db.QueryRow(
// 		"SELECT id, password FROM users WHERE email = $1",
// 		input.Email,
// 	).Scan(&userID, &storedPassword)

// 	if err != nil || !checkPasswordHash(input.Password, storedPassword) {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	token, err := generateToken(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"user_id": userID,
// 		"token":   token,
// 		"message": "Login successful",
// 	})
// }


package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tejaswini22199/task-management-system/authservice/routes"
	database "github.com/tejaswini22199/task-management-system/database"
)

func main() {
	log.Println("Starting Auth Service...")

	// Initialize Database Connection
	database.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterAuthRoutes(r)

	log.Println("Auth Service is running on port 8000")
	r.Run(":8000")
}
