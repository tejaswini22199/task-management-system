package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tejaswini22199/task-management-system/authservice/middleware"
	"github.com/tejaswini22199/task-management-system/taskservice/controllers"
)

func RegisterTaskRoutes(router *gin.Engine) {
	taskRoutes := router.Group("/tasks")
	fmt.Println("line 13")
	taskRoutes.Use(middleware.AuthMiddleware())
	fmt.Println("line 15")

	taskRoutes.POST("", controllers.CreateTask)
	taskRoutes.GET("", controllers.GetTasks)
	taskRoutes.GET("/:id", controllers.GetTaskByID)
	taskRoutes.PUT("/:id", controllers.UpdateTask)
	taskRoutes.GET("/status/:status", controllers.GetTasksByStatus)
	taskRoutes.DELETE("/:id", controllers.DeleteTask)
}
