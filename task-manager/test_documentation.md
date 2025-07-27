# Test Documentation

## Overview

This document provides comprehensive documentation for the test suite of the Task Manager API, which follows Clean Architecture principles. The test suite covers unit tests, integration tests, repository tests, and comprehensive end-to-end tests across different layers of the application.

## Test Structure

The test suite is organized into the following directories and files:

```
task-manager/tests/
├── domain/
│   └── domain_test.go          # Domain entity tests
├── repositories/
│   ├── task_repository_test.go # Task repository integration tests
│   └── user_repository_test.go # User repository integration tests
├── usecases/
    ├── task_usecases_test.go   # Task use case unit tests
    └── user_usecases_test.go   # User use case unit tests
└── End-to-End_test.go          # Comprehensive E2E API tests
```

## Test Configuration

### Makefile Commands

The project includes a `Makefile` with the following test-related commands:

```makefile
.PHONY: test
test:
	go test -v ./tests/... -cover

.PHONY: test-coverage
test-coverage:
	go test -v ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
```

- `make test`: Runs all tests with verbose output and coverage information
- `make test-coverage`: Generates detailed coverage reports in HTML format

### Test Dependencies

The test suite uses the following testing frameworks and libraries:

- **testify/suite**: For structured test suites with setup/teardown
- **testify/assert**: For assertions and test validation
- **MongoDB Go Driver**: For database integration tests
- **godotenv**: For loading environment variables in tests

## Domain Layer Tests

### File: `tests/domain/domain_test.go`

Tests the core domain entities to ensure proper structure and field assignment.

#### Test Cases:

1. **TestTask**: Validates Task entity structure
   - Tests all field assignments (ID, Title, Description, DueDate, Status)
   - Ensures primitive.ObjectID handling works correctly

2. **TestUser**: Validates User entity structure
   - Tests all field assignments (ID, Username, Password, Role)
   - Ensures proper data type handling

#### Coverage:
- ✅ Task entity validation
- ✅ User entity validation
- ✅ MongoDB ObjectID integration

## Repository Layer Tests

### File: `tests/repositories/task_repository_test.go`

Integration tests for the Task Repository that interact with a real MongoDB instance.

#### Test Setup:
- Uses `TestMain` for database connection management
- Connects to MongoDB using environment variables
- Creates isolated test database (`test_taskdb`)
- Implements proper cleanup after tests

#### Test Suite: `TaskRepoTestSuite`

**Setup Methods:**
- `SetupSuite()`: Initializes database and collection
- `SetupTest()`: Cleans collection before each test

**Test Cases:**

1. **TestTaskCreation**
   - Creates a new task in the database
   - Validates ID generation and field persistence
   - Ensures proper error handling

2. **TestGetTaskByID**
   - **Positive Case**: Retrieves existing task by ID
   - **Negative Case**: Returns error for non-existent task
   - Tests ObjectID conversion and validation

3. **TestRetrieveAllTasks**
   - Inserts multiple tasks
   - Validates retrieval of all tasks
   - Tests collection iteration

4. **TestUpdateTask**
   - Creates initial task
   - Updates task fields
   - Validates changes persistence
   - Ensures ID remains unchanged

5. **TestRemoveTask**
   - Creates task for deletion
   - Performs deletion operation
   - Validates task no longer exists

#### Coverage:
- ✅ CRUD operations
- ✅ Error handling
- ✅ Database integration
- ✅ ObjectID validation

### File: `tests/repositories/user_repository_test.go`

Integration tests for the User Repository with mocked dependencies.

#### Mock Implementations:

1. **MockJWTService**
   - `GenerateToken()`: Returns dummy token
   - `ValidateToken()`: Returns dummy claims

2. **MockPasswordService**
   - `HashPassword()`: Returns dummy hash
   - `ComparePassword()`: Always returns success

#### Test Suite: `AuthRepoTestSuite`

**Setup Methods:**
- Creates unique index on username field
- Initializes repository with mock services
- Cleans database before each test

**Test Cases:**

1. **TestUserRegistration**
   - **Success Case**: Registers new user successfully
   - **Duplicate Case**: Rejects duplicate usernames
   - Tests username uniqueness constraint

2. **TestFindUserByUsername**
   - **Success Case**: Locates user by username
   - **Not Found Case**: Returns error for missing user
   - Tests database query functionality

#### Coverage:
- ✅ User registration
- ✅ Username uniqueness
- ✅ User lookup
- ✅ Error handling
- ✅ Mock service integration

## End-to-End Tests

### File: `End-to-End_test.go`

Comprehensive end-to-end tests that validate the entire application workflow from HTTP requests to database persistence. These tests use a real MongoDB instance and test the complete API functionality.

#### Test Setup:
- Uses `TestMain` for MongoDB connection management
- Creates isolated test database (`e2e_test_taskdb`)
- Implements proper cleanup after all tests
- Uses Gin test mode for HTTP testing
- Real HTTP request/response testing with `httptest`

#### Test Suite: `E2ETestSuite`

**Setup Methods:**
- `SetupSuite()`: Initializes database, collections, and router
- `TearDownSuite()`: Cleans up database and connections
- `SetupTest()`: Cleans collections before each test
- `setupUsersForTaskTests()`: Helper to create test users

**Helper Methods:**
- `makeRequest()`: Creates and executes HTTP requests
- `parseResponse()`: Parses JSON responses into structs

#### Test Cases:

**1. TestCompleteUserAuthenticationFlow**
- **User Registration**: Tests admin and regular user registration
- **Role Assignment**: Validates first user becomes admin, subsequent users become regular users
- **User Login**: Tests successful authentication and token generation
- **Duplicate Prevention**: Ensures username uniqueness
- **Invalid Credentials**: Tests rejection of wrong passwords

**Coverage:**
- ✅ Complete registration workflow
- ✅ Role-based user creation
- ✅ JWT token generation
- ✅ Authentication validation
- ✅ Error handling for duplicates and invalid credentials

**2. TestCompleteTaskManagementFlow**
- **Task Creation**: Admin creates tasks successfully
- **Authorization**: Regular users cannot create tasks
- **Task Retrieval**: Get all tasks and individual tasks by ID
- **Task Updates**: Admin can modify existing tasks
- **Task Deletion**: Admin can remove tasks
- **Permission Validation**: Regular users cannot modify/delete tasks

**Coverage:**
- ✅ Full CRUD operations for tasks
- ✅ Role-based authorization
- ✅ HTTP status code validation
- ✅ Data persistence verification
- ✅ Error handling for unauthorized access

**3. TestUserManagementFlow**
- **User Lookup**: Get user information by username
- **User Promotion**: Admin promotes regular users to admin role
- **Authorization**: Regular users cannot promote others
- **Error Handling**: Proper responses for non-existent users

**Coverage:**
- ✅ User information retrieval
- ✅ Role promotion functionality
- ✅ Admin-only operations
- ✅ Not found error handling

**4. TestAuthenticationAndAuthorizationEdgeCases**
- **Missing Tokens**: Access protected endpoints without authentication
- **Invalid Tokens**: Use malformed or expired tokens
- **Role Restrictions**: Regular users accessing admin endpoints
- **Malformed Requests**: Invalid JSON and empty required fields
- **Invalid IDs**: Malformed ObjectID formats

**Coverage:**
- ✅ Authentication edge cases
- ✅ Authorization boundary testing
- ✅ Input validation
- ✅ Error response consistency

**5. TestCompleteApplicationWorkflow**
- **Full Workflow**: Complete user journey from registration to task management
- **Multi-User Scenarios**: Multiple users with different roles
- **State Transitions**: User promotion and role changes
- **Data Consistency**: Verify operations across multiple entities

**Coverage:**
- ✅ End-to-end application workflow
- ✅ Multi-user interactions
- ✅ State management
- ✅ Cross-entity operations

**6. TestPerformanceAndStress**
- **Concurrent Operations**: Multiple simultaneous task creation
- **Rapid Requests**: High-frequency authentication requests
- **Load Testing**: System behavior under concurrent load
- **Timeout Handling**: Proper timeout management

**Coverage:**
- ✅ Concurrent request handling
- ✅ Performance under load
- ✅ System stability
- ✅ Resource management

**7. TestDataValidationAndEdgeCases**
- **Data Types**: Various input data formats and types
- **Field Validation**: Empty, long, and special character inputs
- **Unicode Support**: International characters and emojis
- **Boundary Testing**: Maximum length inputs

**Coverage:**
- ✅ Input data validation
- ✅ Unicode and special character handling
- ✅ Boundary condition testing
- ✅ Data type flexibility

**8. TestDatabaseStateConsistency**
- **State Verification**: Direct database validation after operations
- **CRUD Consistency**: Ensure API operations match database state
- **Data Integrity**: Verify data persistence and updates
- **Cleanup Verification**: Confirm deletions are properly executed

**Coverage:**
- ✅ Database state consistency
- ✅ Data integrity validation
- ✅ Operation persistence
- ✅ Cleanup verification

#### Key Features:

**Real Database Integration:**
- Uses actual MongoDB instance for realistic testing
- Isolated test database to prevent data contamination
- Proper connection management and cleanup

**HTTP Testing:**
- Real HTTP request/response cycle testing
- Proper header handling (Authorization, Content-Type)
- Status code validation
- JSON request/response parsing

**Authentication Testing:**
- JWT token generation and validation
- Role-based access control testing
- Session management across requests

**Concurrent Testing:**
- Goroutine-based concurrent operations
- Channel-based result collection
- Timeout handling for concurrent operations

**Data Validation:**
- Direct database state verification
- Cross-reference API responses with database content
- Consistency checks across operations

## Use Case Layer Tests

### File: `tests/usecases/task_usecases_test.go`

Unit tests for Task Use Cases using repository mocks.

#### Mock Implementation: `StubTaskRepo`

Implements `domain.TaskRepository` interface with configurable behavior:
- `OnCreate`: Mock for CreateTask
- `OnFind`: Mock for GetTaskByID
- `OnFetch`: Mock for GetAllTasks
- `OnUpdate`: Mock for UpdateTask
- `OnRemove`: Mock for DeleteTask

#### Test Suite: `TaskUseCaseSuite`

**Test Cases:**

1. **TestCreateTask**
   - Validates task creation flow
   - Tests ID generation
   - Ensures proper data flow from use case to repository

2. **TestGetTaskByID**
   - Tests task retrieval by ID
   - Validates ID parameter passing
   - Ensures proper data return

3. **TestGetAllTasks**
   - Tests retrieval of multiple tasks
   - Validates data collection handling
   - Tests empty and populated collections

4. **TestUpdateTask**
   - Tests task modification
   - Validates ID preservation
   - Tests field updates

5. **TestDeleteTask**
   - Tests task removal
   - Validates ID parameter passing
   - Tests successful deletion

#### Coverage:
- ✅ All CRUD operations
- ✅ Data flow validation
- ✅ Repository interaction
- ✅ Error propagation

### File: `tests/usecases/user_usecases_test.go`

Unit tests for User Use Cases using repository mocks.

#### Mock Implementation: `StubRepo`

Implements `domain.UserRepository` interface:
- `OnRegister`: Mock for RegisterUser
- `OnLogin`: Mock for LoginUser
- `OnPromote`: Mock for PromoteUser
- `OnFindByUsername`: Mock for GetUserByUsername

#### Test Suite: `UserUseCaseSuite`

**Test Cases:**

1. **TestRegisterUser**
   - **Success Case**: Successful user registration
   - **Failure Case**: Username already taken
   - Tests role assignment and ID generation

2. **TestLoginUser**
   - **Success Case**: Valid credentials authentication
   - **Failure Case**: Invalid credentials rejection
   - Tests token generation flow

3. **TestPromoteUser**
   - **Success Case**: User promotion to admin
   - **Failure Case**: User not found error
   - Tests role modification

4. **TestGetUserByUsername**
   - **Success Case**: User found by username
   - **Failure Case**: User not found error
   - Tests user lookup functionality

#### Coverage:
- ✅ User registration flow
- ✅ Authentication process
- ✅ User promotion
- ✅ User lookup
- ✅ Error handling

## Test Execution

### Running Tests

1. **All Tests with Coverage**:
   ```bash
   make test
   ```

2. **Generate Coverage Report**:
   ```bash
   make test-coverage
   ```
   This generates `coverage.html` for detailed coverage analysis.

3. **Run Specific Test Suite**:
   ```bash
   go test -v ./tests/domain/...
   go test -v ./tests/repositories/...
   go test -v ./tests/usecases/...
   go test -v ./End-to-End_test.go
   ```

4. **Run Only E2E Tests**:
   ```bash
   go test -v ./End-to-End_test.go
   ```

### Environment Setup

Tests require the following environment variables:

- `MONGODB_URI`: MongoDB connection string for integration tests
- Environment variables should be defined in `.env` file

### Test Database

- Integration tests use separate test databases (`test_taskdb` for unit/integration tests, `e2e_test_taskdb` for E2E tests)
- Database is automatically cleaned up after test execution
- Each test starts with a clean state

### E2E Test Requirements

- **MongoDB Instance**: Running MongoDB server (local or remote)
- **Environment Variables**: Proper `.env` configuration
- **Network Access**: Ability to connect to MongoDB
- **Clean State**: Tests create and clean their own database

## Coverage Analysis

Based on the coverage file, the current test coverage includes:

### Covered Components:
- ✅ Domain entities (Task, User)
- ✅ Repository implementations (Task, User)
- ✅ Use case implementations (Task, User)
- ✅ Infrastructure services (JWT, Password, Auth Middleware)
- ✅ Controllers and routing
- ✅ Complete API endpoints (E2E)
- ✅ Authentication and authorization flows
- ✅ Database integration and consistency
- ✅ Error handling and edge cases
- ✅ Concurrent operations and performance

### Coverage Statistics:
- Repository layer: Comprehensive coverage of CRUD operations
- Use case layer: Full business logic coverage
- Infrastructure: Service implementations covered
- Controllers: HTTP handler coverage
- E2E: Complete API workflow coverage
- Integration: Database and HTTP integration coverage

### Test Types Coverage:
- **Unit Tests**: ✅ Use cases with mocked dependencies
- **Integration Tests**: ✅ Repository layer with real database
- **End-to-End Tests**: ✅ Complete API workflows with real HTTP and database
- **Performance Tests**: ✅ Concurrent operations and load testing
- **Security Tests**: ✅ Authentication and authorization validation
## Best Practices Implemented

1. **Test Isolation**: Each test runs in isolation with clean state
2. **Mock Usage**: External dependencies are properly mocked
3. **Integration Testing**: Real database interactions tested
4. **Error Testing**: Both success and failure scenarios covered
5. **Setup/Teardown**: Proper test lifecycle management
6. **Structured Testing**: Using testify/suite for organized tests
7. **End-to-End Validation**: Complete workflow testing
8. **Concurrent Testing**: Multi-threaded operation validation
9. **Database Consistency**: Direct database state verification
10. **HTTP Protocol Testing**: Real request/response cycle validation
11. **Security Testing**: Authentication and authorization validation
12. **Performance Testing**: Load and stress testing capabilities

## Test Execution Strategies

### Development Testing
```bash
# Quick unit tests during development
go test -v ./tests/usecases/...

# Integration tests for database changes
go test -v ./tests/repositories/...

# Full E2E validation before commits
go test -v ./End-to-End_test.go
```

### CI/CD Pipeline
```bash
# Complete test suite with coverage
make test-coverage

# E2E tests in staging environment
MONGODB_URI=staging_uri go test -v ./End-to-End_test.go
```

### Performance Testing
```bash
# Run specific performance tests
go test -v ./End-to-End_test.go -run TestPerformanceAndStress
```

## Test Metrics and Reporting

### Coverage Metrics
- **Line Coverage**: Measures code execution during tests
- **Branch Coverage**: Validates all conditional paths
- **Function Coverage**: Ensures all functions are tested

### Performance Metrics
- **Response Time**: API endpoint response times
- **Throughput**: Requests per second under load
- **Concurrency**: Simultaneous operation handling
- **Resource Usage**: Memory and CPU utilization during tests

### Quality Metrics
- **Test Success Rate**: Percentage of passing tests
- **Error Coverage**: Percentage of error scenarios tested
- **Edge Case Coverage**: Boundary condition testing completeness
## Conclusion

The test suite provides comprehensive coverage of the Task Manager API's core functionality across all layers of the application. It follows testing best practices with proper separation between unit tests (use cases), integration tests (repositories), and end-to-end tests (complete API workflows). 

The addition of comprehensive E2E tests ensures that:
- **Complete Workflows** are validated from HTTP request to database persistence
- **Real-World Scenarios** are tested with actual HTTP and database interactions
- **Performance Characteristics** are validated under concurrent load
- **Security Measures** are properly tested across all endpoints
- **Data Consistency** is maintained across all operations
- **Error Handling** is comprehensive and user-friendly

The test suite now covers the full spectrum of testing needs:
- **Unit Tests**: Fast, isolated business logic validation
- **Integration Tests**: Database and service integration validation  
- **End-to-End Tests**: Complete application workflow validation
- **Performance Tests**: Load and concurrency validation
- **Security Tests**: Authentication and authorization validation

This comprehensive testing approach ensures high confidence in the application's reliability, performance, and security across all deployment environments.