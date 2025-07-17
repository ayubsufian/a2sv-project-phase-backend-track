# Task Manager API

## Overview

A simple, in-memory REST API for task management built with Go and Gin.

## Features

- **Full CRUD**:
  - List all tasks
  - Retrieve a specific task
  - Create a task
  - Update a task
  - Delete a task
- **Validation**
- **Concurrency-safe**

## Installation

1. Ensure **Go 1.20+** is installed
2. Clone the repo

```bash
git clone https://github.com/ayubsufian/a2sv-project-phase-backend-track.git
cd task_manager
```

3. Install dependencies

```bash
go mod tidy
```

## Running the API

```bash
go run main.go
```

The server listens on http://localhost:8080.

```

```

## API Endpoints

[API Documentation](https://documenter.getpostman.com/view/46809956/2sB34imLWd)

## Validation Errors

Example response format for validation failure:

```json
{
  "errors": {
    "Title": "required",
    "DueDate": "duedate",
    "Status": "oneof"
  }
}
```

## Testing

curl or Postman can be used for testing the API.

Example curl Commands

```bash

# Create

curl -X POST http://localhost:8080/tasks \
 -H "Content-Type: application/json" \
 -d '{"title":"Test","due_date":"2025-07-30T14:00:00Z","status":"pending"}'

# List

curl http://localhost:8080/tasks

# Update

curl -X PUT http://localhost:8080/tasks/1 \
 -H "Content-Type: application/json" \
 -d '{"title":"Updated","due_date":"2025-08-01T10:00:00Z","status":"completed"}'

# Delete

curl -X DELETE http://localhost:8080/tasks/1
```
