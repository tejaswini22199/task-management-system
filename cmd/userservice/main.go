package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/your-username/task-management-system/internal/user/handlers"
    "github.com/your-username/task-management-system/pkg/database"
)

func main() {
    // Initialize DB
    db, err := database.InitDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize router
    r := gin.Default()
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(db)

    // Routes
    r.POST("/auth/register", authHandler.Register)
    r.POST("/auth/login", authHandler.Login)
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "User service is healthy"})
    })

    port := os.Getenv("SERVICE_PORT")
    if port == "" {
        port = "8081"
    }
    r.Run(":" + port)
}