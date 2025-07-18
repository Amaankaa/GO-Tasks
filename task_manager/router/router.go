package router

import (
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	// r.POST("/promote", controllers.PromoteUser) // To be implemented

	tasks := r.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.GET("", controllers.GetTasks)
		tasks.GET(":id", controllers.GetTaskByID)
		tasks.POST("", middleware.AdminOnly(), controllers.CreateTask)
		tasks.PUT(":id", middleware.AdminOnly(), controllers.UpdateTask)
		tasks.DELETE(":id", middleware.AdminOnly(), controllers.DeleteTask)
	}

	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		users.POST(":id/promote", controllers.Promote)
	}

	return r
}
