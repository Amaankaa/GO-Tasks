# Task Management API Documentation

## Endpoints

### GET /tasks
- **Description:** Get a list of all tasks.
- **Response:**
```
Status: 200 OK
[
  {
    "id": 1,
    "title": "Task Title",
    "description": "Task Description",
    "due_date": "2024-06-01",
    "status": "pending"
  },
  ...
]
```

### GET /tasks/:id
- **Description:** Get details of a specific task.
- **Response:**
```
Status: 200 OK
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "due_date": "2024-06-01",
  "status": "pending"
}
```
- **Error:**
```
Status: 404 Not Found
{
  "error": "Task not found"
}
```

### POST /tasks
- **Description:** Create a new task.
- **Request:**
```
{
  "title": "Task Title",
  "description": "Task Description",
  "due_date": "2024-06-01",
  "status": "pending"
}
```
- **Response:**
```
Status: 201 Created
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "due_date": "2024-06-01",
  "status": "pending"
}
```
- **Error:**
```
Status: 400 Bad Request
{
  "error": "<error message>"
}
```

### PUT /tasks/:id
- **Description:** Update a specific task.
- **Request:**
```
{
  "title": "Updated Title",
  "description": "Updated Description",
  "due_date": "2024-06-02",
  "status": "completed"
}
```
- **Response:**
```
Status: 200 OK
{
  "id": 1,
  "title": "Updated Title",
  "description": "Updated Description",
  "due_date": "2024-06-02",
  "status": "completed"
}
```
- **Error:**
```
Status: 400 Bad Request | 404 Not Found
{
  "error": "<error message>"
}
```

### DELETE /tasks/:id
- **Description:** Delete a specific task.
- **Response:**
```
Status: 204 No Content
```
- **Error:**
```
Status: 404 Not Found
{
  "error": "Task not found"
}
``` 