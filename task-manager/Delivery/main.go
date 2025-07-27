package main

import (
	"context"
	"log"
	"os"
	"time"

	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	infrastructure "task-manager/Infrastructure"
	"task-manager/Repositories"
	"task-manager/Usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI from environment variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set in environment")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Select database and collections
	db := client.Database("taskdb")
	taskCollection := db.Collection("tasks")
	userCollection := db.Collection("users")

	// Initialize services
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService()

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(taskCollection)
	userRepo := repositories.NewUserRepository(userCollection, jwtService, passwordService)

	// Initialize usecases
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
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
