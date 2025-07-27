# Task Manager API Documentation

## Overview

The Task Manager API is a RESTful web service built with Go using the Gin framework and MongoDB for data persistence. The application follows Clean Architecture principles with clear separation of concerns across different layers. It provides comprehensive task management functionality with user authentication and role-based authorization.

## Architecture

The application is structured using Clean Architecture with the following layers:

### 1. Domain Layer (`Domain/`)
Contains core business entities and interfaces, independent of external frameworks.

### 2. Use Cases Layer (`Usecases/`)
Implements application-specific business rules and orchestrates data flow.

### 3. Infrastructure Layer (`Infrastructure/`)
Handles external dependencies like JWT tokens, password hashing, and authentication middleware.

### 4. Repository Layer (`Repositories/`)
Abstracts data access logic and implements database operations.

### 5. Delivery Layer (`Delivery/`)
Handles HTTP requests/responses and contains the application entry point.

## Technology Stack

- **Language**: Go 1.20
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Environment Management**: godotenv

## Project Structure

```
task-manager/
├── Delivery/
│   ├── main.go                 # Application entry point
│   ├── controllers/
│   │   └── controller.go       # HTTP request handlers
│   └── routers/
│       └── router.go           # Route definitions and middleware setup
├── Domain/
│   └── domain.go               # Core entities and interfaces
├── Infrastructure/
│   ├── auth_middleWare.go      # Authentication middleware
│   ├── jwt_service.go          # JWT token operations
│   └── password_service.go     # Password hashing operations
├── Repositories/
│   ├── task_repository.go      # Task data access layer
│   └── user_repository.go      # User data access layer
├── Usecases/
│   ├── task_usecases.go        # Task business logic
│   └── user_usecases.go        # User business logic
└── tests/                      # Test suite
```

## Data Models

### Task Entity

```go
type Task struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    DueDate     string             `bson:"due_date" json:"due_date"`
    Status      string             `bson:"status" json:"status"`
}
```

**Fields:**
- `ID`: Unique MongoDB ObjectID
- `Title`: Task title (required)
- `Description`: Detailed task description
- `DueDate`: Due date in string format
- `Status`: Current task status (e.g., "pending", "completed", "in-progress")

### User Entity

```go
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username string             `bson:"username" json:"username"`
    Password string             `bson:"password" json:"password"`
    Role     string             `bson:"role" json:"role"`
}
```

**Fields:**
- `ID`: Unique MongoDB ObjectID
- `Username`: Unique username (required)
- `Password`: Hashed password (required)
- `Role`: User role ("user" or "admin")

### Login Response

```go
type LoginResponse struct {
    ID       primitive.ObjectID `json:"id"`
    Username string             `json:"username"`
    Token    string             `json:"token"`
}
```

## Authentication & Authorization

### JWT Authentication
The application uses JWT tokens for authentication with the following claims:
- `_id`: User ID
- `username`: Username
- `role`: User role

### Authorization Levels
1. **Public**: No authentication required
2. **Authenticated**: Valid JWT token required
3. **Admin Only**: Admin role required

### Middleware
- `AuthMiddleware()`: Validates JWT tokens and extracts user information
- `AdminOnly()`: Restricts access to admin users only

## API Endpoints

### Base URL
```
http://localhost:8080
```

### Authentication Endpoints

#### 1. User Registration
**POST** `/register`

Registers a new user in the system. The first registered user automatically becomes an admin.

**Request Body:**
```json
{
    "username": "string",
    "password": "string"
}
```

**Response (201 Created):**
```json
{
    "id": "ObjectID",
    "username": "string",
    "role": "user|admin",
    "password": ""
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or username already taken
```json
{
    "error": "username already taken"
}
```

**Business Logic:**
- Validates username and password are not empty
- Checks for username uniqueness
- Hashes password using bcrypt
- Assigns "admin" role to first user, "user" role to subsequent users
- Returns user data with password field cleared

---

#### 2. User Login
**POST** `/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
    "username": "string",
    "password": "string"
}
```

**Response (200 OK):**
```json
{
    "id": "ObjectID",
    "username": "string",
    "token": "JWT_TOKEN_STRING"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid credentials
```json
{
    "error": "invalid username or password"
}
```

**Business Logic:**
- Validates username is provided
- Looks up user by username
- Compares provided password with stored hash
- Generates JWT token with user claims
- Returns user info and token

---

### Task Management Endpoints

All task endpoints require authentication via `Authorization: Bearer <token>` header.

#### 3. Get All Tasks
**GET** `/tasks`

Retrieves all tasks in the system.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK):**
```json
[
    {
        "id": "ObjectID",
        "title": "string",
        "description": "string",
        "due_date": "string",
        "status": "string"
    }
]
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Database error

**Business Logic:**
- Requires valid authentication
- Retrieves all tasks from database
- Returns empty array if no tasks exist

---

#### 4. Get Task by ID
**GET** `/tasks/:id`

Retrieves a specific task by its ID.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `id`: MongoDB ObjectID of the task

**Response (200 OK):**
```json
{
    "id": "ObjectID",
    "title": "string",
    "description": "string",
    "due_date": "string",
    "status": "string"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Task not found
```json
{
    "error": "Task not found"
}
```

**Business Logic:**
- Validates ObjectID format
- Looks up task in database
- Returns task data if found

---

#### 5. Create Task
**POST** `/tasks`

Creates a new task. **Admin access required.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "title": "string",
    "description": "string",
    "due_date": "string",
    "status": "string"
}
```

**Response (201 Created):**
```json
{
    "id": "ObjectID",
    "title": "string",
    "description": "string",
    "due_date": "string",
    "status": "string"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
```json
{
    "error": "Admin access required"
}
```
- `500 Internal Server Error`: Database error

**Business Logic:**
- Requires admin role
- Validates request body
- Generates new ObjectID
- Saves task to database
- Returns created task with ID

---

#### 6. Update Task
**PUT** `/tasks/:id`

Updates an existing task. **Admin access required.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `id`: MongoDB ObjectID of the task

**Request Body:**
```json
{
    "title": "string",
    "description": "string",
    "due_date": "string",
    "status": "string"
}
```

**Response (200 OK):**
```json
{
    "id": "ObjectID",
    "title": "string",
    "description": "string",
    "due_date": "string",
    "status": "string"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or ID format
```json
{
    "Error": "error_message"
}
```
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: Task not found
```json
{
    "message": "Task not Found"
}
```

**Business Logic:**
- Requires admin role
- Validates ObjectID format
- Updates specified fields in database
- Returns updated task data

---

#### 7. Delete Task
**DELETE** `/tasks/:id`

Deletes a task. **Admin access required.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `id`: MongoDB ObjectID of the task

**Response (204 No Content):**
No response body.

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: Task not found
```json
{
    "error": "Task not found"
}
```

**Business Logic:**
- Requires admin role
- Validates ObjectID format
- Removes task from database
- Returns 204 status on success

---

### User Management Endpoints

#### 8. Get User by Username
**GET** `/users/:username`

Retrieves user information by username. **Authentication required.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `username`: Username to look up

**Response (200 OK):**
```json
{
    "id": "ObjectID",
    "username": "string",
    "password": "hashed_password",
    "role": "string"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found
```json
{
    "error": "User not found"
}
```
- `500 Internal Server Error`: Database error

**Business Logic:**
- Requires authentication
- Looks up user by username
- Returns user data including hashed password

---

#### 9. Promote User
**POST** `/users/:id/promote`

Promotes a user to admin role. **Admin access required.**

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `id`: MongoDB ObjectID of the user

**Response (200 OK):**
```json
{
    "id": "ObjectID",
    "username": "string",
    "password": "",
    "role": "admin"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID
```json
{
    "error": "invalid user ID"
}
```
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: User not found
```json
{
    "error": "user not found"
}
```
- `500 Internal Server Error`: Database error

**Business Logic:**
- Requires admin role
- Validates ObjectID format
- Updates user role to "admin"
- Returns updated user data with password cleared

---

## Error Handling

### Standard Error Response Format
```json
{
    "error": "error_message"
}
```

### HTTP Status Codes Used
- `200 OK`: Successful GET/PUT requests
- `201 Created`: Successful POST requests
- `204 No Content`: Successful DELETE requests
- `400 Bad Request`: Invalid request data or parameters
- `401 Unauthorized`: Authentication required or invalid
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side errors

## Security Features

### Password Security
- Passwords are hashed using bcrypt with default cost
- Original passwords are never stored or returned in responses
- Password comparison uses constant-time comparison

### JWT Security
- Tokens include user ID, username, and role claims
- Tokens are validated on each protected request
- Uses HMAC-SHA256 signing method
- Development secret key (should be replaced in production)

### Authorization
- Role-based access control (RBAC)
- Admin-only endpoints for task management
- User context available in request handlers

## Database Operations

### MongoDB Collections
- `tasks`: Stores task documents
- `users`: Stores user documents with unique username index

### Connection Management
- Connection timeout: 10 seconds
- Operation timeout: 5 seconds
- Automatic connection cleanup on application shutdown

### Data Validation
- ObjectID format validation
- Required field validation
- Unique constraint on usernames

## Environment Configuration

### Required Environment Variables
- `MONGODB_URI`: MongoDB connection string (defaults to `mongodb://localhost:27017`)

### Default Configuration
- Server port: `:8080`
- Database name: `taskdb`
- JWT secret: `your_dev_secret_key` (development only)

## Running the Application

### Prerequisites
- Go 1.20 or higher
- MongoDB instance running
- Environment variables configured

### Installation & Startup
1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set up environment variables in `.env` file
4. Run the application: `go run Delivery/main.go`

### Server Output
```
Server starting on :8080
```

## API Usage Examples

### Register First User (Becomes Admin)
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Create Task (Admin Only)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Complete project",
    "description": "Finish the task manager API",
    "due_date": "2024-12-31",
    "status": "pending"
  }'
```

### Get All Tasks
```bash
curl -X GET http://localhost:8080/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Development Notes

### Code Organization
- Clean Architecture principles followed
- Dependency injection used throughout
- Interface-based design for testability
- Separation of concerns maintained

### Error Handling Patterns
- Repository layer returns domain-specific errors
- Use case layer handles business logic errors
- Controller layer converts to appropriate HTTP responses
- Consistent error message formatting

### Future Enhancements
- Add task assignment to specific users
- Implement task categories and tags
- Add task priority levels
- Implement task due date notifications
- Add pagination for task lists
- Implement task search and filtering
- Add audit logging
- Implement refresh token mechanism
- Add rate limiting
- Implement CORS support

## Conclusion

The Task Manager API provides a robust foundation for task management with proper authentication, authorization, and clean architecture. The API is designed to be scalable, maintainable, and follows Go best practices for web service development.