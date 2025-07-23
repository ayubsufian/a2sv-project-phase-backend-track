# Task Manager API

## Overview

A secure RESTful API in Go + Gin with MongoDB, offering full CRUD for tasks (title, optional description, mandatory due_date, status: pending/completed), JWT-based authentication, user roles (user/admin), input validation, and protected endpoints under /api/tasks.

## Prerequisites

- Go 1.20+
- MongoDB Atlas or local MongoDB instance
- .env file with MongoDB connection URI and JWT_SECRET.

## Features

- **Full CRUD**:
  - List all tasks
  - Retrieve a specific task
  - Create a task
  - Update a task
  - Delete a task
- **Validation**
- **User auth**:
  - POST /register: create account (username, password; role field).
  - POST /login: authenticate and retrieve a 24-hour JWT token.
- **Protected endpoints**:
  - All POST, PUT, DELETE, and GET /api/tasks require Authorization: Bearer <token>
  - Admin-only access under /api/admin/dashboard

## Setup and MongoDB Configuration

1. Create .env in project root:

```bash
MONGODB_URI="mongodb+srv://<user>:<pass>@cluster0.<xyz>.mongodb.net/taskdb?retryWrites=true&w=majority"
JWT_SECRET="your_jwt_secret"
```

2. Load env var using github.com/joho/godotenv in main.go, then read MONGODB_URI

## Installation

1. Ensure **Go 1.20+** is installed
2. Clone the repo

```bash
git clone https://github.com/ayubsufian/a2sv-project-phase-backend-track.git
cd task_manager_auth
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

[API Documentation](https://documenter.getpostman.com/view/46809956/2sB34oAbyP)

## Authentication Flow

### Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
 -H "Content-Type: application/json" \
 -d '{"username":"alice","password":"pass123"}'
```

Response:

```json
{ "token": "<JWT_TOKEN>" }
```

## Working with Tasks (Protected Endpoints)

Use the returned JWT as:

```bash
Authorization: Bearer <JWT_TOKEN>
```

### Create Task

```bash
curl -X POST http://localhost:8080/api/tasks \
 -H "Authorization: Bearer $TOKEN" \
 -H "Content-Type: application/json" \
 -d '{"title":"Buy groceries","description":"Milk, eggs","due_date":"2025-08-01T12:00:00Z","status":"pending"}'
```

### List Tasks

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/tasks
```

### Update Task

```bash
curl -X PUT http://localhost:8080/api/tasks/<id> \
 -H "Authorization: Bearer $TOKEN" \
 -H "Content-Type: application/json" \
 -d '{"title":"Updated","due_date":"2025-08-01T10:00:00Z","status":"completed"}'
```

### Delete Task

```bash
curl -X DELETE http://localhost:8080/api/tasks/<id> \
 -H "Authorization: Bearer $TOKEN"
```

## Admin Endpoint

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/api/admin/dashboard
```

Returns:

```json
{ "message": "Welcome Admin" }
```

403 Forbidden for non-admin users.
