package domain_test

import (
	"testing"
	"task-manager/Domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTask(t *testing.T) {
	id := primitive.NewObjectID()
	task := domain.Task{
		ID:          id,
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     "2024-12-31",
		Status:      "pending",
	}

	assert.Equal(t, id, task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Test Description", task.Description)
	assert.Equal(t, "2024-12-31", task.DueDate)
	assert.Equal(t, "pending", task.Status)
}

func TestUser(t *testing.T) {
	id := primitive.NewObjectID()
	user := domain.User{
		ID:       id,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}

	assert.Equal(t, id, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "user", user.Role)
}