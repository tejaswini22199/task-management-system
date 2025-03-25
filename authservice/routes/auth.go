package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tejaswini22199/task-management-system/authservice/controllers"
)

func RegisterAuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", controllers.RegisterUser)
		authRoutes.POST("/login", controllers.LoginUser)
	}
}
