package main

import (
	"log"
	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	"task-manager/Infrastructure"
	"task-manager/Repositories"
	"task-manager/Usecases"
)

func main() {
	// Initialize services
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService()

	// Initialize repositories
	taskRepo, err := repositories.NewTaskRepository()
	if err != nil {
		log.Fatalf("Failed to initialize task repository: %v", err)
	}

	userRepo, err := repositories.NewUserRepository(jwtService, passwordService)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	// Initialize use cases
	taskUsecase := usecases.NewTaskUsecase(taskRepo)
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Initialize controllers
	controller := controllers.NewController(taskUsecase, userUsecase)

	// Initialize middleware
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	// Setup router
	r := routers.SetupRouter(controller, authMiddleware)

	// Start server
	log.Println("Server starting on :8080")
	r.Run(":8080")
}