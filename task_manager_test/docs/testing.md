# Task Manager API - Testing Documentation

## 1. Introduction

This document outlines the testing strategy, setup, and execution procedures for the Task Manager API. The project employs a multi-layered testing approach, leveraging both **unit tests** and **integration tests** to ensure code quality, correctness, and reliability.

The primary testing library used is **Testify**, which provides a structured suite-based approach (`testify/suite`), powerful assertion tools (`testify/assert`), and a complete mocking toolkit (`testify/mock`).

- **Unit Tests:** These are used to test individual components (controllers, use cases, services) in complete isolation. Dependencies are mocked to ensure that tests are fast, repeatable, and target specific logic.
- **Integration Tests:** These are used to test components that interact with external infrastructure. In this project, they specifically validate that the repository layer functions correctly against a real MongoDB database instance, which is managed by Docker.

## 2. Prerequisites

To run the full test suite locally, you will need the following software installed and configured:

- **Go:** Version 1.18 or higher.
- **Docker Desktop:** Required to run the MongoDB instance for integration tests. Ensure the Docker daemon is running.

## 3. Running the Tests Locally

All tests can be run from the root directory of the project (`task_manager_test`).

### 3.1. Running Unit Tests Only

Unit tests are fast and have no external dependencies. They are ideal for running frequently during development.

From a terminal in the project root, run the following command:

```bash
go test -v ./internal/delivery/... ./internal/usecase/... ./internal/service/...
```

- `-v`: Enables verbose mode to see the output of each test as it runs.
- `./...`: The `...` wildcard tells Go to run tests in the specified package and all its sub-packages.

### 3.2. Running Integration Tests Only

Integration tests require a running MongoDB instance. Follow these steps precisely in your **PowerShell** terminal.

**Step 1: Start a MongoDB Docker Container**

This command starts a temporary database named `mongo-test-db` and makes it available on `localhost:27017`.

```powershell
docker run --name mongo-test-db -p 27017:27017 -d mongo
```

**Step 2: Set the Environment Variable**

The test suite reads the database connection string from an environment variable.

```powershell
$env:MONGODB_URI_TEST = "mongodb://localhost:27017"
```

_(Note: This variable is only set for the current terminal session.)_

**Step 3: Run the Repository Tests**

```powershell
go test -v ./internal/repository/...
```

**Step 4: Clean Up**

After the tests are complete, stop and remove the container to free up resources.

````powershell
docker stop mongo-test-db
docker rm mongo-test-db```

### 3.3. Running the Full Test Suite

To run every single test (unit and integration), follow the integration test setup and then run `go test` on the entire project.

```powershell
# 1. Start Docker container
docker run --name mongo-test-db -p 27017:27017 -d mongo

# 2. Set Environment Variable
$env:MONGODB_URI_TEST = "mongodb://localhost:27017"

# 3. Run all tests
go test -v ./...

# 4. Clean up Docker container
docker stop mongo-test-db
docker rm mongo-test-db
````

## 4. Test Suite Overview

### 4.1. Delivery Layer (`internal/delivery`)

- **`task_controller_test.go` & `user_controller_test.go`**

  - **Type:** Unit Test
  - **Methodology:** Mocks the `TaskUsecase` and `UserUsecase` dependencies. Uses `httptest` to create in-memory HTTP requests and record responses.
  - **Verifies:** Correct mapping of HTTP requests to controller methods, proper JSON request binding, correct HTTP status codes for all success and error paths (e.g., `200`, `201`, `400`, `401`, `404`, `409`), and accurate JSON response bodies.

- **`auth_test.go`**

  - **Type:** Unit Test
  - **Methodology:** Mocks the `IJWTService` dependency. Tests the `AuthMiddleware` and `AdminOnly` middleware functions by injecting them into a test router.
  - **Verifies:** Rejection of requests with no token, malformed tokens, or invalid tokens. Correctly populates the request context with user claims on success. Correctly enforces role-based access for the `AdminOnly` middleware.

- **`router_test.go`**

  - **Type:** Unit Test
  - **Methodology:** Injects mock controllers and services into `SetupRouter`. Iterates through the router's registered routes to verify their properties.
  - **Verifies:** That all required API routes are registered with the correct HTTP method and path. Verifies that middleware (`AuthMiddleware`, `AdminOnly`) is correctly applied to the appropriate route groups.

### 4.2. Usecase Layer (`internal/usecase`)

- **`task_usecase_test.go` & `user_usecase_test.go`**

  - **Type:** Unit Test
  - **Methodology:** Mocks all repository (`ITaskRepository`, `IUserRepository`) and service (`IPasswordService`, `IJWTService`) dependencies.
  - **Verifies:** Core business logic in isolation. Tests the sequence and flow of operations (e.g., in `Register`, it verifies that `Hash` is called before `Create`). Validates business rules (e.g., task title cannot be empty) and correct error handling and propagation from dependencies.

### 4.3. Service Layer (`internal/service`)

- **`password_service_test.go` & `jwt_service_test.go`**

  - **Type:** Unit Test
  - **Methodology:** These services have no external dependencies, so they are tested directly without mocks.
  - **Verifies:** `password_service` correctly hashes passwords and compares hashes. `jwt_service` correctly generates and validates tokens, including failure cases for expired tokens, invalid signatures, and malformed tokens.

### 4.4. Repository Layer (`internal/repository`)

- **`task_repository_test.go` & `user_repository_test.go`**

  - **Type:** Integration Test
  - **Methodology:** Connects to a real MongoDB instance running in a Docker container. The test suite manages the database connection and ensures a clean state for each test by creating fresh collections and indexes in `SetupTest` and dropping them in `TearDownTest`.
  - **Verifies:** Correct CRUD operations against the database. Verifies that database constraints (like unique indexes) are working as expected and that database errors (like `ErrNoDocuments` or duplicate key errors) are correctly mapped to application-specific errors (like `usecase.ErrNotFound` or `usecase.ErrUserAlreadyExists`).

## 5. Test Coverage

High test coverage is a key objective to ensure code quality.

### 5.1. Generating a Coverage Report

To generate a test coverage report for the entire project, run the following command from the root directory:

```bash
go test -cover ./...
```

### 5.2. Generating a Visual Coverage Report

For a more detailed, visual view of which lines of code are covered, you can generate an HTML report.

**Step 1: Generate a coverage profile**

```bash
go test -coverprofile=coverage.out ./...
```

**Step 2: View the HTML report in your browser**

```bash
go tool cover -html=coverage.out
```

This command will open a new tab in your browser with an interactive report, highlighting covered (green), uncovered (red), and untested code.

## 6. Issues Encountered During Testing

During the development of this test suite, several key issues were identified and resolved. Documenting them serves as a reference for future development.

1.  **Issue: `UserRepository` Duplicate User Test Failed**

    - **Symptom:** `TestCreate_Fails_When_UserAlreadyExists` was failing. The test expected `usecase.ErrUserAlreadyExists`, but the repository `Create` method returned `nil`, indicating a successful insertion of a duplicate user.
    - **Root Cause:** The MongoDB test collection was missing a unique index on the `username` field. Furthermore, the index was initially created in `SetupSuite`, but the `TearDownTest` function dropped the entire collection, which also removed the index.
    - **Resolution:** The unique index creation logic was moved from `SetupSuite` to `SetupTest`. This ensures the index is correctly configured before every single test run, guaranteeing test isolation and correctness.

2.  **Issue: `AuthMiddleware` Test Panicked on Type Conversion**

    - **Symptom:** `TestAuthMiddleware_Success` panicked with the error `interface conversion: interface {} is map[string]interface {}, not jwt.MapClaims`.
    - **Root Cause:** The `IJWTService` mock was configured to return a generic `map[string]interface{}`. The mock's auto-generated code, however, strictly expected the `jwt.MapClaims` type as defined in the interface.
    - **Resolution:** The test was corrected to create the claims map using the specific `jwt.MapClaims` type, satisfying the mock's type assertion.

3.  **Issue: Docker Network Timeout on Image Pull**

    - **Symptom:** The `docker run -d mongo` command failed with a `TLS handshake timeout` error.
    - **Root Cause:** This was an external environmental issue related to local networking (e.g., an unstable connection, firewall, or VPN) preventing Docker from connecting to Docker Hub.
    - **Resolution:** Resolved by stabilizing the network connection and ensuring no firewalls or VPNs were interfering with Docker's network access.
