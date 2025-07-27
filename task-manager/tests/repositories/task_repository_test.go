package test_repositories

import (
	"context"
	"log"
	"os"
	"testing"

	domain "task-manager/Domain"
	repositories "task-manager/Repositories"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testMongoClient *mongo.Client

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Unable to load environment config")
	}

	connStr := os.Getenv("MONGODB_URI")
	if connStr == "" {
		log.Fatal("MONGODB_URI not found in environment")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("MongoDB ping unsuccessful: %v", err)
	}

	log.Println("✔ Connected to test MongoDB")
	testMongoClient = client

	exitCode := m.Run()

	testDb := client.Database("test_db")
	if err := testDb.Drop(context.Background()); err != nil {
		log.Printf("⚠ Could not drop test DB: %v", err)
	}

	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("⚠ MongoDB disconnect error: %v", err)
	}

	os.Exit(exitCode)
}

type TaskRepoTestSuite struct {
	suite.Suite
	repo domain.TaskRepository
	coll *mongo.Collection
	db   *mongo.Database
}

func TestTaskRepoIntegration(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("MongoDB not initialized, skipping tests.")
	}
	suite.Run(t, new(TaskRepoTestSuite))
}

func (suite *TaskRepoTestSuite) SetupSuite() {
	suite.db = testMongoClient.Database("test_taskdb")
	suite.coll = suite.db.Collection("tasks")
}

func (suite *TaskRepoTestSuite) SetupTest() {
	_, err := suite.coll.DeleteMany(context.Background(), bson.D{})
	suite.Require().NoError(err, "Database cleanup failed")
	suite.repo = repositories.NewTaskRepository(suite.coll)
}

func (suite *TaskRepoTestSuite) TestTaskCreation() {
	input := &domain.Task{
		Title:       "Integration Task",
		Description: "Testing task creation flow",
		DueDate:     "July 17, 2024",
		Status:      "pending",
	}

	result, err := suite.repo.CreateTask(*input)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.False(result.ID.IsZero(), "Task ID should be generated")
	suite.Equal(input.Title, result.Title)
}

func (suite *TaskRepoTestSuite) TestGetTaskByID() {
	suite.Run("Should return task if present", func() {
		task := &domain.Task{
			ID:    primitive.NewObjectID(),
			Title: "Lookup Task",
		}
		_, err := suite.coll.InsertOne(context.Background(), task)
		suite.Require().NoError(err)

		retrieved, err := suite.repo.GetTaskByID(task.ID.Hex())

		suite.Require().NoError(err)
		suite.Require().NotNil(retrieved)
		suite.Equal(task.ID, retrieved.ID)
	})

	suite.Run("Should return error if task is missing", func() {
		unknownID := primitive.NewObjectID().Hex()
		_, err := suite.repo.GetTaskByID(unknownID)

		suite.Require().Error(err)
	})
}

func (suite *TaskRepoTestSuite) TestRetrieveAllTasks() {
	docs := []interface{}{
		&domain.Task{ID: primitive.NewObjectID(), Title: "Alpha"},
		&domain.Task{ID: primitive.NewObjectID(), Title: "Beta"},
	}
	_, err := suite.coll.InsertMany(context.Background(), docs)
	suite.Require().NoError(err)

	tasks, err := suite.repo.GetAllTasks()
	suite.Require().NoError(err)
	suite.Len(tasks, 2)
}

func (suite *TaskRepoTestSuite) TestUpdateTask() {
	existing := &domain.Task{
		ID:     primitive.NewObjectID(),
		Title:  "Initial Title",
		Status: "pending",
	}
	_, err := suite.coll.InsertOne(context.Background(), existing)
	suite.Require().NoError(err)

	patch := &domain.Task{
		Title:  "Final Title",
		Status: "completed",
	}
	updated, err := suite.repo.UpdateTask(existing.ID.Hex(), *patch)

	suite.Require().NoError(err)
	suite.Require().NotNil(updated)
	suite.Equal("Final Title", updated.Title)
	suite.Equal(existing.ID, updated.ID)
}

func (suite *TaskRepoTestSuite) TestRemoveTask() {
	task := &domain.Task{
		ID:    primitive.NewObjectID(),
		Title: "Obsolete Task",
	}
	_, err := suite.coll.InsertOne(context.Background(), task)
	suite.Require().NoError(err)

	err = suite.repo.DeleteTask(task.ID.Hex())
	suite.Require().NoError(err)

	_, err = suite.repo.GetTaskByID(task.ID.Hex())
	suite.Error(err, "Expected error after deletion")
}
