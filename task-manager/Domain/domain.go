package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task entity
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     string             `bson:"due_date" json:"due_date"`
	Status      string             `bson:"status" json:"status"`
}

// User represents a user entity
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Role     string             `bson:"role" json:"role"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Token    string             `json:"token"`
}

// TaskRepository interface defines task data access operations
type TaskRepository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id string) (Task, error)
	CreateTask(task Task) (Task, error)
	UpdateTask(id string, task Task) (Task, error)
	DeleteTask(id string) error
}

// UserRepository interface defines user data access operations
type UserRepository interface {
	RegisterUser(user User) (User, error)
	LoginUser(user User) (LoginResponse, error)
	PromoteUser(id string) (User, error)
	GetUserByUsername(username string) (User, error)
}

// TaskUsecase interface defines task business logic operations
type TaskUsecase interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id string) (Task, error)
	CreateTask(task Task) (Task, error)
	UpdateTask(id string, task Task) (Task, error)
	DeleteTask(id string) error
}

// UserUsecase interface defines user business logic operations
type UserUsecase interface {
	RegisterUser(user User) (User, error)
	LoginUser(user User) (LoginResponse, error)
	PromoteUser(id string) (User, error)
}

// JWTService interface defines JWT operations
type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

// PasswordService interface defines password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}