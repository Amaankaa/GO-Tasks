# Task Manager API - Clean Architecture

This is a refactored version of the Task Manager API following Clean Architecture principles with proper separation of concerns.

## Architecture Overview

The application follows Clean Architecture with the following layers:

### Domain Layer (`Domain/`)
- Contains core business entities (Task, User)
- Defines interfaces for repositories and use cases
- No dependencies on external frameworks

### Use Cases Layer (`Usecases/`)
- Contains application-specific business rules
- Implements use case interfaces defined in Domain
- Orchestrates data flow between repositories and controllers

### Infrastructure Layer (`Infrastructure/`)
- Implements external dependencies and services
- JWT service for token generation and validation
- Password service for hashing and comparison
- Authentication middleware

### Repository Layer (`Repositories/`)
- Abstracts data access logic
- Implements repository interfaces defined in Domain
- Handles MongoDB operations

### Delivery Layer (`Delivery/`)
- Handles HTTP requests and responses
- Controllers invoke appropriate use case methods
- Router configuration and server setup

## Folder Structure

```
task-manager/
├── Delivery/
│   ├── main.go                 # Application entry point
│   ├── controllers/
│   │   └── controller.go       # HTTP request handlers
│   └── routers/
│       └── router.go           # Route definitions
├── Domain/
│   └── domain.go               # Core entities and interfaces
├── Infrastructure/
│   ├── auth_middleWare.go      # Authentication middleware
│   ├── jwt_service.go          # JWT token operations
│   └── password_service.go     # Password hashing operations
├── Repositories/
│   ├── task_repository.go      # Task data access
│   └── user_repository.go      # User data access
└── Usecases/
    ├── task_usecases.go        # Task business logic
    └── user_usecases.go        # User business logic
```

## Key Improvements

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Dependency Inversion**: High-level modules don't depend on low-level modules
3. **Interface Segregation**: Small, focused interfaces
4. **Testability**: Easy to mock dependencies for unit testing
5. **Maintainability**: Changes in one layer don't affect others

## Running the Application

1. Ensure MongoDB is running (default: `mongodb://localhost:27017`)
2. Navigate to the project directory
3. Run: `go run Delivery/main.go`
4. Server starts on `:8080`

## API Endpoints

- `POST /register` - Register a new user
- `POST /login` - User login
- `GET /tasks` - Get all tasks (authenticated)
- `GET /tasks/:id` - Get task by ID (authenticated)
- `POST /tasks` - Create task (admin only)
- `PUT /tasks/:id` - Update task (admin only)
- `DELETE /tasks/:id` - Delete task (admin only)
- `POST /users/:id/promote` - Promote user to admin (admin only)

## Environment Variables

- `MONGODB_URI` - MongoDB connection string (optional, defaults to localhost)

The application maintains the same functionality as the original while providing better code organization and maintainability.