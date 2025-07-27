package routers

import (
	"task-manager/Delivery/controllers"
	"task-manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *controllers.Controller, authMiddleware *infrastructure.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	// Protected task routes
	tasks := r.Group("/tasks")
	tasks.Use(authMiddleware.AuthMiddleware())
	{
		tasks.GET("", controller.GetTasks)
		tasks.GET(":id", controller.GetTaskByID)
		tasks.POST("", authMiddleware.AdminOnly(), controller.CreateTask)
		tasks.PUT(":id", authMiddleware.AdminOnly(), controller.UpdateTask)
		tasks.DELETE(":id", authMiddleware.AdminOnly(), controller.DeleteTask)
	}

	// Protected user routes
	users := r.Group("/users")
	users.Use(authMiddleware.AuthMiddleware())
	{
		users.GET(":username", controller.GetUserByUsername)
	}
	users.Use(authMiddleware.AuthMiddleware(), authMiddleware.AdminOnly())
	{
		users.POST(":id/promote", controller.Promote)
	}

	return r
}