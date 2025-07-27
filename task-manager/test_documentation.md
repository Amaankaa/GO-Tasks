# Test Documentation

## Overview

This document provides comprehensive documentation for the test suite of the Task Manager API, which follows Clean Architecture principles. The test suite covers unit tests, integration tests, and repository tests across different layers of the application.

## Test Structure

The test suite is organized into the following directories:

```
task-manager/tests/
├── domain/
│   └── domain_test.go          # Domain entity tests
├── repositories/
│   ├── task_repository_test.go # Task repository integration tests
│   └── user_repository_test.go # User repository integration tests
└── usecases/
    ├── task_usecases_test.go   # Task use case unit tests
    └── user_usecases_test.go   # User use case unit tests
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
   ```

### Environment Setup

Tests require the following environment variables:

- `MONGODB_URI`: MongoDB connection string for integration tests
- Environment variables should be defined in `.env` file

### Test Database

- Integration tests use a separate test database (`test_taskdb`)
- Database is automatically cleaned up after test execution
- Each test starts with a clean state

## Coverage Analysis

Based on the coverage file, the current test coverage includes:

### Covered Components:
- ✅ Domain entities (Task, User)
- ✅ Repository implementations (Task, User)
- ✅ Use case implementations (Task, User)
- ✅ Infrastructure services (JWT, Password, Auth Middleware)
- ✅ Controllers and routing

### Coverage Statistics:
- Repository layer: Comprehensive coverage of CRUD operations
- Use case layer: Full business logic coverage
- Infrastructure: Service implementations covered
- Controllers: HTTP handler coverage

## Best Practices Implemented

1. **Test Isolation**: Each test runs in isolation with clean state
2. **Mock Usage**: External dependencies are properly mocked
3. **Integration Testing**: Real database interactions tested
4. **Error Testing**: Both success and failure scenarios covered
5. **Setup/Teardown**: Proper test lifecycle management
6. **Structured Testing**: Using testify/suite for organized tests

## Recommendations

1. **Add Controller Tests**: Direct HTTP endpoint testing
2. **Add Middleware Tests**: Authentication and authorization testing
3. **Performance Tests**: Load testing for database operations
4. **End-to-End Tests**: Full API workflow testing
5. **Test Data Builders**: Create helper functions for test data generation

## Conclusion

The test suite provides comprehensive coverage of the Task Manager API's core functionality. It follows testing best practices with proper separation between unit tests (use cases) and integration tests (repositories). The use of mocks and test suites ensures maintainable and reliable tests that validate both business logic and data persistence layers.