# Task Manager API

## Overview

A RESTful API in Go + Gin with MongoDB, offering full CRUD for tasks featuring title, optional description, mandatory future due_date, and status (pending/completed).

## Prerequisites

- Go 1.20+
- MongoDB Atlas or local MongoDB instance
- .env file with MongoDB connection URI

## Features

- **Full CRUD**:
  - List all tasks
  - Retrieve a specific task
  - Create a task
  - Update a task
  - Delete a task
- **Validation**

## Setup and MongoDB Configuration

1. Create .env in project root:

```bash
MONGODB_URI="mongodb+srv://<user>:<pass>@cluster0.<xyz>.mongodb.net/taskdb?retryWrites=true&w=majority"
```

2. Load env var using github.com/joho/godotenv in main.go, then read MONGODB_URI
3. Connect to MongoDB in main.go with a 10-second timeout, ping health check, and collection initialization:

```go
client, err := mongo.Connect(ctx, options.Client().
    ApplyURI(uri).
    SetServerSelectionTimeout(15*time.Second).
    SetRetryWrites(true))
```

4. Retry logic: up to 5 attempts with exponential backoff.
5. Global init:

```go
data.Client = client
data.TasksCollection = client.Database("taskdb").Collection("tasks")
```

## Installation

1. Ensure **Go 1.20+** is installed
2. Clone the repo

```bash
git clone https://github.com/ayubsufian/a2sv-project-phase-backend-track.git
cd task_manager_mongodb
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

## API Endpoints

[API Documentation](https://documenter.getpostman.com/view/46809956/2sB34kEeko)

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

curl -X PUT http://localhost:8080/tasks/<id> \
 -H "Content-Type: application/json" \
 -d '{"title":"Updated","due_date":"2025-08-01T10:00:00Z","status":"completed"}'

# Delete

curl -X DELETE http://localhost:8080/tasks/<id>
```
