package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	domain "task-manager/Domain"
	infrastructure "task-manager/Infrastructure"
	repositories "task-manager/Repositories"
	usecases "task-manager/Usecases"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// E2ETestSuite represents the end-to-end test suite
type E2ETestSuite struct {
	suite.Suite
	router       *gin.Engine
	client       *mongo.Client
	db           *mongo.Database
	taskColl     *mongo.Collection
	userColl     *mongo.Collection
	adminToken   string
	userToken    string
	adminUserID  string
	regularUserID string
	testTaskID   string
}

// TestE2ETestSuite runs the end-to-end test suite
func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

// SetupSuite initializes the test environment
func (suite *E2ETestSuite) SetupSuite() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Get MongoDB URI
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	suite.Require().NoError(err, "Failed to connect to MongoDB")

	// Ping MongoDB
	err = client.Ping(ctx, nil)
	suite.Require().NoError(err, "Failed to ping MongoDB")

	suite.client = client
	suite.db = client.Database("e2e_test_taskdb")
	suite.taskColl = suite.db.Collection("tasks")
	suite.userColl = suite.db.Collection("users")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize services
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService()

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(suite.taskColl)
	userRepo := repositories.NewUserRepository(suite.userColl, jwtService, passwordService)

	// Initialize use cases
	taskUsecase := usecases.NewTaskUsecase(taskRepo)
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Initialize controllers
	controller := controllers.NewController(taskUsecase, userUsecase)

	// Initialize middleware
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	// Setup router
	suite.router = routers.SetupRouter(controller, authMiddleware)

	log.Println("✅ E2E Test Suite initialized successfully")
}

// TearDownSuite cleans up the test environment
func (suite *E2ETestSuite) TearDownSuite() {
	if suite.client != nil {
		// Drop test database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := suite.db.Drop(ctx)
		if err != nil {
			log.Printf("Warning: Could not drop test database: %v", err)
		}

		// Disconnect from MongoDB
		err = suite.client.Disconnect(ctx)
		if err != nil {
			log.Printf("Warning: Could not disconnect from MongoDB: %v", err)
		}
	}
	log.Println("✅ E2E Test Suite cleaned up successfully")
}

// SetupTest cleans the database before each test
func (suite *E2ETestSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean collections
	_, err := suite.taskColl.DeleteMany(ctx, bson.D{})
	suite.Require().NoError(err)

	_, err = suite.userColl.DeleteMany(ctx, bson.D{})
	suite.Require().NoError(err)

	// Reset tokens and IDs
	suite.adminToken = ""
	suite.userToken = ""
	suite.adminUserID = ""
	suite.regularUserID = ""
	suite.testTaskID = ""
}

// Helper method to make HTTP requests
func (suite *E2ETestSuite) makeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		suite.Require().NoError(err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, path, reqBody)
	suite.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	return w
}

// Helper method to parse JSON response
func (suite *E2ETestSuite) parseResponse(w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	suite.Require().NoError(err, "Failed to parse response: %s", w.Body.String())
}

// Test 1: Complete User Registration and Authentication Flow
func (suite *E2ETestSuite) TestCompleteUserAuthenticationFlow() {
	suite.Run("Register first user as admin", func() {
		adminUser := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		w := suite.makeRequest("POST", "/register", adminUser, "")
		suite.Equal(http.StatusCreated, w.Code)

		var response domain.User
		suite.parseResponse(w, &response)

		suite.Equal("admin", response.Username)
		suite.Equal("admin", response.Role)
		suite.NotEmpty(response.ID)
		suite.Empty(response.Password) // Password should be cleared
		suite.adminUserID = response.ID.Hex()
	})

	suite.Run("Register second user as regular user", func() {
		regularUser := map[string]string{
			"username": "user",
			"password": "user123",
		}

		w := suite.makeRequest("POST", "/register", regularUser, "")
		suite.Equal(http.StatusCreated, w.Code)

		var response domain.User
		suite.parseResponse(w, &response)

		suite.Equal("user", response.Username)
		suite.Equal("user", response.Role)
		suite.NotEmpty(response.ID)
		suite.regularUserID = response.ID.Hex()
	})

	suite.Run("Login admin user", func() {
		loginData := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		w := suite.makeRequest("POST", "/login", loginData, "")
		suite.Equal(http.StatusOK, w.Code)

		var response domain.LoginResponse
		suite.parseResponse(w, &response)

		suite.Equal("admin", response.Username)
		suite.NotEmpty(response.Token)
		suite.adminToken = response.Token
	})

	suite.Run("Login regular user", func() {
		loginData := map[string]string{
			"username": "user",
			"password": "user123",
		}

		w := suite.makeRequest("POST", "/login", loginData, "")
		suite.Equal(http.StatusOK, w.Code)

		var response domain.LoginResponse
		suite.parseResponse(w, &response)

		suite.Equal("user", response.Username)
		suite.NotEmpty(response.Token)
		suite.userToken = response.Token
	})

	suite.Run("Reject duplicate username registration", func() {
		duplicateUser := map[string]string{
			"username": "admin",
			"password": "different123",
		}

		w := suite.makeRequest("POST", "/register", duplicateUser, "")
		suite.Equal(http.StatusBadRequest, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "username already taken")
	})

	suite.Run("Reject invalid login credentials", func() {
		invalidLogin := map[string]string{
			"username": "admin",
			"password": "wrongpassword",
		}

		w := suite.makeRequest("POST", "/login", invalidLogin, "")
		suite.Equal(http.StatusUnauthorized, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "invalid username or password")
	})
}

// Test 2: Complete Task Management Flow
func (suite *E2ETestSuite) TestCompleteTaskManagementFlow() {
	// Setup users first
	suite.setupUsersForTaskTests()

	suite.Run("Admin creates a task", func() {
		newTask := map[string]string{
			"title":       "Complete E2E Tests",
			"description": "Implement comprehensive end-to-end testing",
			"due_date":    "2024-12-31",
			"status":      "pending",
		}

		w := suite.makeRequest("POST", "/tasks", newTask, suite.adminToken)
		suite.Equal(http.StatusCreated, w.Code)

		var response domain.Task
		suite.parseResponse(w, &response)

		suite.Equal("Complete E2E Tests", response.Title)
		suite.Equal("Implement comprehensive end-to-end testing", response.Description)
		suite.Equal("2024-12-31", response.DueDate)
		suite.Equal("pending", response.Status)
		suite.NotEmpty(response.ID)
		suite.testTaskID = response.ID.Hex()
	})

	suite.Run("Regular user cannot create task", func() {
		newTask := map[string]string{
			"title":       "Unauthorized Task",
			"description": "This should fail",
			"status":      "pending",
		}

		w := suite.makeRequest("POST", "/tasks", newTask, suite.userToken)
		suite.Equal(http.StatusForbidden, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "Admin access required")
	})

	suite.Run("Get all tasks (authenticated user)", func() {
		w := suite.makeRequest("GET", "/tasks", nil, suite.userToken)
		suite.Equal(http.StatusOK, w.Code)

		var response []domain.Task
		suite.parseResponse(w, &response)

		suite.Len(response, 1)
		suite.Equal("Complete E2E Tests", response[0].Title)
	})

	suite.Run("Get task by ID", func() {
		path := fmt.Sprintf("/tasks/%s", suite.testTaskID)
		w := suite.makeRequest("GET", path, nil, suite.userToken)
		suite.Equal(http.StatusOK, w.Code)

		var response domain.Task
		suite.parseResponse(w, &response)

		suite.Equal("Complete E2E Tests", response.Title)
		suite.Equal(suite.testTaskID, response.ID.Hex())
	})

	suite.Run("Admin updates task", func() {
		updatedTask := map[string]string{
			"title":       "Complete E2E Tests - Updated",
			"description": "Implement comprehensive end-to-end testing with full coverage",
			"due_date":    "2024-12-25",
			"status":      "in-progress",
		}

		path := fmt.Sprintf("/tasks/%s", suite.testTaskID)
		w := suite.makeRequest("PUT", path, updatedTask, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var response domain.Task
		suite.parseResponse(w, &response)

		suite.Equal("Complete E2E Tests - Updated", response.Title)
		suite.Equal("Implement comprehensive end-to-end testing with full coverage", response.Description)
		suite.Equal("2024-12-25", response.DueDate)
		suite.Equal("in-progress", response.Status)
		suite.Equal(suite.testTaskID, response.ID.Hex())
	})

	suite.Run("Regular user cannot update task", func() {
		updatedTask := map[string]string{
			"title":  "Unauthorized Update",
			"status": "completed",
		}

		path := fmt.Sprintf("/tasks/%s", suite.testTaskID)
		w := suite.makeRequest("PUT", path, updatedTask, suite.userToken)
		suite.Equal(http.StatusForbidden, w.Code)
	})

	suite.Run("Get non-existent task returns 404", func() {
		w := suite.makeRequest("GET", "/tasks/507f1f77bcf86cd799439011", nil, suite.userToken)
		suite.Equal(http.StatusNotFound, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "Task not found")
	})

	suite.Run("Admin deletes task", func() {
		path := fmt.Sprintf("/tasks/%s", suite.testTaskID)
		w := suite.makeRequest("DELETE", path, nil, suite.adminToken)
		suite.Equal(http.StatusNoContent, w.Code)
		suite.Empty(w.Body.String())
	})

	suite.Run("Verify task is deleted", func() {
		path := fmt.Sprintf("/tasks/%s", suite.testTaskID)
		w := suite.makeRequest("GET", path, nil, suite.userToken)
		suite.Equal(http.StatusNotFound, w.Code)
	})

	suite.Run("Regular user cannot delete task", func() {
		// Create a new task first
		newTask := map[string]string{
			"title":  "Task to Delete",
			"status": "pending",
		}

		w := suite.makeRequest("POST", "/tasks", newTask, suite.adminToken)
		suite.Equal(http.StatusCreated, w.Code)

		var task domain.Task
		suite.parseResponse(w, &task)

		// Try to delete with regular user
		path := fmt.Sprintf("/tasks/%s", task.ID.Hex())
		w = suite.makeRequest("DELETE", path, nil, suite.userToken)
		suite.Equal(http.StatusForbidden, w.Code)
	})
}

// Test 3: User Management Flow
func (suite *E2ETestSuite) TestUserManagementFlow() {
	// Setup users first
	suite.setupUsersForTaskTests()

	suite.Run("Get user by username", func() {
		path := "/users/admin"
		w := suite.makeRequest("GET", path, nil, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var response domain.User
		suite.parseResponse(w, &response)

		suite.Equal("admin", response.Username)
		suite.Equal("admin", response.Role)
		suite.NotEmpty(response.Password) // Hashed password should be present
	})

	suite.Run("Get non-existent user returns 404", func() {
		path := "/users/nonexistent"
		w := suite.makeRequest("GET", path, nil, suite.adminToken)
		suite.Equal(http.StatusNotFound, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "User not found")
	})

	suite.Run("Admin promotes regular user", func() {
		path := fmt.Sprintf("/users/%s/promote", suite.regularUserID)
		w := suite.makeRequest("POST", path, nil, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var response domain.User
		suite.parseResponse(w, &response)

		suite.Equal("user", response.Username)
		suite.Equal("admin", response.Role) // Should be promoted to admin
		suite.Empty(response.Password)      // Password should be cleared in response
	})

	suite.Run("Regular user cannot promote others", func() {
		path := fmt.Sprintf("/users/%s/promote", suite.adminUserID)
		w := suite.makeRequest("POST", path, nil, suite.userToken)
		suite.Equal(http.StatusForbidden, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "Admin access required")
	})

	suite.Run("Promote non-existent user returns 404", func() {
		path := "/users/507f1f77bcf86cd799439011/promote"
		w := suite.makeRequest("POST", path, nil, suite.adminToken)
		suite.Equal(http.StatusNotFound, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "user not found")
	})
}

// Test 4: Authentication and Authorization Edge Cases
func (suite *E2ETestSuite) TestAuthenticationAndAuthorizationEdgeCases() {
	// Setup users first
	suite.setupUsersForTaskTests()

	suite.Run("Access protected endpoint without token", func() {
		w := suite.makeRequest("GET", "/tasks", nil, "")
		suite.Equal(http.StatusUnauthorized, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "Authorization header missing or invalid")
	})

	suite.Run("Access protected endpoint with invalid token", func() {
		w := suite.makeRequest("GET", "/tasks", nil, "invalid-token")
		suite.Equal(http.StatusUnauthorized, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "invalid or expired token")
	})

	suite.Run("Access admin endpoint with regular user token", func() {
		newTask := map[string]string{
			"title":  "Unauthorized Task",
			"status": "pending",
		}

		w := suite.makeRequest("POST", "/tasks", newTask, suite.userToken)
		suite.Equal(http.StatusForbidden, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "Admin access required")
	})

	suite.Run("Invalid request body format", func() {
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid-json")))
		suite.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusBadRequest, w.Code)
	})

	suite.Run("Empty required fields in registration", func() {
		emptyUser := map[string]string{
			"username": "",
			"password": "",
		}

		w := suite.makeRequest("POST", "/register", emptyUser, "")
		suite.Equal(http.StatusBadRequest, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "fields cannot be empty")
	})

	suite.Run("Invalid ObjectID format in task operations", func() {
		w := suite.makeRequest("GET", "/tasks/invalid-id", nil, suite.userToken)
		suite.Equal(http.StatusBadRequest, w.Code)

		var errorResponse map[string]string
		suite.parseResponse(w, &errorResponse)
		suite.Contains(errorResponse["error"], "invalid id format")
	})
}

// Test 5: Complete Application Workflow
func (suite *E2ETestSuite) TestCompleteApplicationWorkflow() {
	suite.Run("Complete workflow from registration to task management", func() {
		// Step 1: Register admin user
		adminUser := map[string]string{
			"username": "workflow_admin",
			"password": "admin123",
		}

		w := suite.makeRequest("POST", "/register", adminUser, "")
		suite.Equal(http.StatusCreated, w.Code)

		var adminResponse domain.User
		suite.parseResponse(w, &adminResponse)
		suite.Equal("admin", adminResponse.Role)

		// Step 2: Register regular user
		regularUser := map[string]string{
			"username": "workflow_user",
			"password": "user123",
		}

		w = suite.makeRequest("POST", "/register", regularUser, "")
		suite.Equal(http.StatusCreated, w.Code)

		var userResponse domain.User
		suite.parseResponse(w, &userResponse)
		suite.Equal("user", userResponse.Role)

		// Step 3: Login both users
		w = suite.makeRequest("POST", "/login", adminUser, "")
		suite.Equal(http.StatusOK, w.Code)

		var adminLogin domain.LoginResponse
		suite.parseResponse(w, &adminLogin)
		adminToken := adminLogin.Token

		w = suite.makeRequest("POST", "/login", regularUser, "")
		suite.Equal(http.StatusOK, w.Code)

		var userLogin domain.LoginResponse
		suite.parseResponse(w, &userLogin)
		userToken := userLogin.Token

		// Step 4: Admin creates multiple tasks
		tasks := []map[string]string{
			{
				"title":       "Task 1",
				"description": "First task",
				"status":      "pending",
			},
			{
				"title":       "Task 2",
				"description": "Second task",
				"status":      "in-progress",
			},
			{
				"title":       "Task 3",
				"description": "Third task",
				"status":      "completed",
			},
		}

		var createdTasks []domain.Task
		for _, task := range tasks {
			w = suite.makeRequest("POST", "/tasks", task, adminToken)
			suite.Equal(http.StatusCreated, w.Code)

			var createdTask domain.Task
			suite.parseResponse(w, &createdTask)
			createdTasks = append(createdTasks, createdTask)
		}

		// Step 5: Regular user views all tasks
		w = suite.makeRequest("GET", "/tasks", nil, userToken)
		suite.Equal(http.StatusOK, w.Code)

		var allTasks []domain.Task
		suite.parseResponse(w, &allTasks)
		suite.Len(allTasks, 3)

		// Step 6: Admin updates a task
		updatedTask := map[string]string{
			"title":       "Task 1 - Updated",
			"description": "First task - updated description",
			"status":      "completed",
		}

		path := fmt.Sprintf("/tasks/%s", createdTasks[0].ID.Hex())
		w = suite.makeRequest("PUT", path, updatedTask, adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var updated domain.Task
		suite.parseResponse(w, &updated)
		suite.Equal("Task 1 - Updated", updated.Title)
		suite.Equal("completed", updated.Status)

		// Step 7: Admin promotes regular user
		path = fmt.Sprintf("/users/%s/promote", userResponse.ID.Hex())
		w = suite.makeRequest("POST", path, nil, adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var promotedUser domain.User
		suite.parseResponse(w, &promotedUser)
		suite.Equal("admin", promotedUser.Role)

		// Step 8: Newly promoted admin logs in again to get a new token
		loginData := map[string]string{
			"username": userResponse.Username,
			"password": regularUser["password"], // or whatever password you used for this user
		}
		w = suite.makeRequest("POST", "/login", loginData, "")
		suite.Equal(http.StatusOK, w.Code)

		var promotedLogin domain.LoginResponse
		suite.parseResponse(w, &promotedLogin)
		promotedUserToken := promotedLogin.Token

		// Now use the new token for admin actions
		newAdminTask := map[string]string{
			"title":       "Task by Promoted Admin",
			"description": "Task created by newly promoted admin",
			"status":      "pending",
		}

		w = suite.makeRequest("POST", "/tasks", newAdminTask, promotedUserToken)
		suite.Equal(http.StatusCreated, w.Code)

		// Step 9: Verify final state
		w = suite.makeRequest("GET", "/tasks", nil, adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var finalTasks []domain.Task
		suite.parseResponse(w, &finalTasks)
		suite.Len(finalTasks, 4) // 3 original + 1 new

		// Step 10: Clean up by deleting a task
		path = fmt.Sprintf("/tasks/%s", createdTasks[2].ID.Hex())
		w = suite.makeRequest("DELETE", path, nil, adminToken)
		suite.Equal(http.StatusNoContent, w.Code)

		// Verify deletion
		w = suite.makeRequest("GET", "/tasks", nil, adminToken)
		suite.Equal(http.StatusOK, w.Code)

		suite.parseResponse(w, &finalTasks)
		suite.Len(finalTasks, 3) // Should be 3 after deletion
	})
}

// Helper method to setup users for task tests
func (suite *E2ETestSuite) setupUsersForTaskTests() {
	if suite.adminToken == "" || suite.userToken == "" {
		// Register admin user
		adminUser := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		w := suite.makeRequest("POST", "/register", adminUser, "")
		suite.Require().Equal(http.StatusCreated, w.Code)

		var adminResponse domain.User
		suite.parseResponse(w, &adminResponse)
		suite.adminUserID = adminResponse.ID.Hex()

		// Register regular user
		regularUser := map[string]string{
			"username": "user",
			"password": "user123",
		}

		w = suite.makeRequest("POST", "/register", regularUser, "")
		suite.Require().Equal(http.StatusCreated, w.Code)

		var userResponse domain.User
		suite.parseResponse(w, &userResponse)
		suite.regularUserID = userResponse.ID.Hex()

		// Login admin
		w = suite.makeRequest("POST", "/login", adminUser, "")
		suite.Require().Equal(http.StatusOK, w.Code)

		var adminLogin domain.LoginResponse
		suite.parseResponse(w, &adminLogin)
		suite.adminToken = adminLogin.Token

		// Login regular user
		w = suite.makeRequest("POST", "/login", regularUser, "")
		suite.Require().Equal(http.StatusOK, w.Code)

		var userLogin domain.LoginResponse
		suite.parseResponse(w, &userLogin)
		suite.userToken = userLogin.Token
	}
}

// Test 6: Performance and Stress Testing
func (suite *E2ETestSuite) TestPerformanceAndStress() {
	suite.setupUsersForTaskTests()

	suite.Run("Create multiple tasks concurrently", func() {
		const numTasks = 10
		taskChan := make(chan domain.Task, numTasks)
		errorChan := make(chan error, numTasks)

		// Create tasks concurrently
		for i := 0; i < numTasks; i++ {
			go func(taskNum int) {
				task := map[string]string{
					"title":       fmt.Sprintf("Concurrent Task %d", taskNum),
					"description": fmt.Sprintf("Task created concurrently - %d", taskNum),
					"status":      "pending",
				}

				w := suite.makeRequest("POST", "/tasks", task, suite.adminToken)
				if w.Code != http.StatusCreated {
					errorChan <- fmt.Errorf("failed to create task %d: status %d", taskNum, w.Code)
					return
				}

				var createdTask domain.Task
				err := json.Unmarshal(w.Body.Bytes(), &createdTask)
				if err != nil {
					errorChan <- fmt.Errorf("failed to parse task %d: %v", taskNum, err)
					return
				}

				taskChan <- createdTask
			}(i)
		}

		// Collect results
		var createdTasks []domain.Task
		var errors []error

		for i := 0; i < numTasks; i++ {
			select {
			case task := <-taskChan:
				createdTasks = append(createdTasks, task)
			case err := <-errorChan:
				errors = append(errors, err)
			case <-time.After(10 * time.Second):
				suite.Fail("Timeout waiting for concurrent task creation")
			}
		}

		suite.Empty(errors, "Should have no errors in concurrent task creation")
		suite.Len(createdTasks, numTasks, "Should create all tasks successfully")

		// Verify all tasks exist
		w := suite.makeRequest("GET", "/tasks", nil, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		var allTasks []domain.Task
		suite.parseResponse(w, &allTasks)
		suite.GreaterOrEqual(len(allTasks), numTasks, "Should have at least the created tasks")
	})

	suite.Run("Rapid authentication requests", func() {
		const numRequests = 20
		successChan := make(chan bool, numRequests)
		errorChan := make(chan error, numRequests)

		loginData := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		// Make rapid login requests
		for i := 0; i < numRequests; i++ {
			go func(reqNum int) {
				w := suite.makeRequest("POST", "/login", loginData, "")
				if w.Code != http.StatusOK {
					errorChan <- fmt.Errorf("login request %d failed: status %d", reqNum, w.Code)
					return
				}

				var response domain.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					errorChan <- fmt.Errorf("failed to parse login response %d: %v", reqNum, err)
					return
				}

				if response.Token == "" {
					errorChan <- fmt.Errorf("empty token in response %d", reqNum)
					return
				}

				successChan <- true
			}(i)
		}

		// Collect results
		var successes int
		var errors []error

		for i := 0; i < numRequests; i++ {
			select {
			case <-successChan:
				successes++
			case err := <-errorChan:
				errors = append(errors, err)
			case <-time.After(15 * time.Second):
				suite.Fail("Timeout waiting for rapid authentication requests")
			}
		}

		suite.Empty(errors, "Should have no errors in rapid authentication")
		suite.Equal(numRequests, successes, "All authentication requests should succeed")
	})
}

// Test 7: Data Validation and Edge Cases
func (suite *E2ETestSuite) TestDataValidationAndEdgeCases() {
	suite.setupUsersForTaskTests()

	suite.Run("Task creation with various data types", func() {
		testCases := []struct {
			name         string
			task         map[string]interface{}
			expectedCode int
		}{
			{
				name: "Valid task with all fields",
				task: map[string]interface{}{
					"title":       "Valid Task",
					"description": "This is a valid task",
					"due_date":    "2024-12-31",
					"status":      "pending",
				},
				expectedCode: http.StatusCreated,
			},
			{
				name: "Task with empty title",
				task: map[string]interface{}{
					"title":       "",
					"description": "Task with empty title",
					"status":      "pending",
				},
				expectedCode: http.StatusCreated, // API doesn't validate empty title
			},
			{
				name: "Task with very long title",
				task: map[string]interface{}{
					"title":       string(make([]byte, 1000)), // Very long title
					"description": "Task with long title",
					"status":      "pending",
				},
				expectedCode: http.StatusCreated,
			},
			{
				name: "Task with special characters",
				task: map[string]interface{}{
					"title":       "Task with special chars: !@#$%^&*()",
					"description": "Description with unicode: 你好世界 🌍",
					"status":      "pending",
				},
				expectedCode: http.StatusCreated,
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				w := suite.makeRequest("POST", "/tasks", tc.task, suite.adminToken)
				suite.Equal(tc.expectedCode, w.Code, "Test case: %s", tc.name)

				if tc.expectedCode == http.StatusCreated {
					var response domain.Task
					suite.parseResponse(w, &response)
					suite.NotEmpty(response.ID)
				}
			})
		}
	})

	suite.Run("User registration edge cases", func() {
		testCases := []struct {
			name         string
			user         map[string]interface{}
			expectedCode int
		}{
			{
				name: "Username with special characters",
				user: map[string]interface{}{
					"username": "user@domain.com",
					"password": "password123",
				},
				expectedCode: http.StatusCreated,
			},
			{
				name: "Very long username",
				user: map[string]interface{}{
					"username": string(make([]byte, 100)),
					"password": "password123",
				},
				expectedCode: http.StatusCreated,
			},
			{
				name: "Username with unicode",
				user: map[string]interface{}{
					"username": "用户名",
					"password": "password123",
				},
				expectedCode: http.StatusCreated,
			},
			{
				name: "Very long password",
				user: map[string]interface{}{
					"username": "longpassuser",
					"password": string(make([]byte, 200)),
				},
				expectedCode: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				w := suite.makeRequest("POST", "/register", tc.user, "")
				suite.Equal(tc.expectedCode, w.Code, "Test case: %s", tc.name)

				if tc.expectedCode == http.StatusCreated {
					var response domain.User
					suite.parseResponse(w, &response)
					suite.NotEmpty(response.ID)
				}
			})
		}
	})
}

// Test 8: Database State Consistency
func (suite *E2ETestSuite) TestDatabaseStateConsistency() {
	suite.setupUsersForTaskTests()

	suite.Run("Verify database state after operations", func() {
		// Create a task
		newTask := map[string]string{
			"title":       "Consistency Test Task",
			"description": "Testing database consistency",
			"status":      "pending",
		}

		w := suite.makeRequest("POST", "/tasks", newTask, suite.adminToken)
		suite.Equal(http.StatusCreated, w.Code)

		var createdTask domain.Task
		suite.parseResponse(w, &createdTask)

		// Verify task exists in database directly
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var dbTask domain.Task
		err := suite.taskColl.FindOne(ctx, bson.M{"_id": createdTask.ID}).Decode(&dbTask)
		suite.NoError(err, "Task should exist in database")
		suite.Equal(createdTask.Title, dbTask.Title)

		// Update the task
		updatedTask := map[string]string{
			"title":       "Updated Consistency Test Task",
			"description": "Updated description",
			"status":      "completed",
		}

		path := fmt.Sprintf("/tasks/%s", createdTask.ID.Hex())
		w = suite.makeRequest("PUT", path, updatedTask, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		// Verify update in database
		err = suite.taskColl.FindOne(ctx, bson.M{"_id": createdTask.ID}).Decode(&dbTask)
		suite.NoError(err, "Updated task should exist in database")
		suite.Equal("Updated Consistency Test Task", dbTask.Title)
		suite.Equal("completed", dbTask.Status)

		// Delete the task
		w = suite.makeRequest("DELETE", path, nil, suite.adminToken)
		suite.Equal(http.StatusNoContent, w.Code)

		// Verify deletion in database
		err = suite.taskColl.FindOne(ctx, bson.M{"_id": createdTask.ID}).Decode(&dbTask)
		suite.Error(err, "Deleted task should not exist in database")
		suite.Equal(mongo.ErrNoDocuments, err)
	})

	suite.Run("Verify user state consistency", func() {
		// Check initial user count
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		initialCount, err := suite.userColl.CountDocuments(ctx, bson.M{})
		suite.NoError(err)

		// Register a new user
		newUser := map[string]string{
			"username": "consistency_user",
			"password": "password123",
		}

		w := suite.makeRequest("POST", "/register", newUser, "")
		suite.Equal(http.StatusCreated, w.Code)

		var createdUser domain.User
		suite.parseResponse(w, &createdUser)

		// Verify user count increased
		newCount, err := suite.userColl.CountDocuments(ctx, bson.M{})
		suite.NoError(err)
		suite.Equal(initialCount+1, newCount)

		// Verify user exists in database with hashed password
		var dbUser domain.User
		err = suite.userColl.FindOne(ctx, bson.M{"_id": createdUser.ID}).Decode(&dbUser)
		suite.NoError(err)
		suite.Equal(createdUser.Username, dbUser.Username)
		suite.NotEqual("password123", dbUser.Password) // Should be hashed
		suite.NotEmpty(dbUser.Password)

		// Promote the user
		path := fmt.Sprintf("/users/%s/promote", createdUser.ID.Hex())
		w = suite.makeRequest("POST", path, nil, suite.adminToken)
		suite.Equal(http.StatusOK, w.Code)

		// Verify promotion in database
		err = suite.userColl.FindOne(ctx, bson.M{"_id": createdUser.ID}).Decode(&dbUser)
		suite.NoError(err)
		suite.Equal("admin", dbUser.Role)
	})
}