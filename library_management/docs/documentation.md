# Library Management System Documentation

## Overview

The system is a **console-based (CLI) Library Management System** written in Go. It supports:

- Managing book records (add/remove)
- Registering library members
- Processing book loans and returns
- Viewing available and borrowed book lists

---

## Project Structure

```
library_management/
├── controllers/
│   └── library_controller.go        # CLI input handling and commands
├── docs/
│   └── documentation.md             # This documentation
├── models/
│   └── book.go                      # Book model
│   └── member.go                    # Member model
├── services/
│   └── library_service.go          # Core library logic
└── main.go                          # Application entrypoint
```

---

## Architecture & Components

- **main.go**

  - Initializes the `Library` service and launches the CLI.

- **controllers/**

  - Contains CLI handling logic (`StartCLI`) that interprets user commands and calls service methods.

- **services/**

  - Contains `Library` implementation: handles book and member data, borrowing and returning logic.

- **models/**

  - Defines core data structures: `Book` and `Member`.

- **docs/**

  - Contains system documentation.

---

## Core Concepts

### Models

- **Book**

  - Fields: `ID` (int), `Title` (string), `Author` (string), `Status` (string — “Available” or “Borrowed”)

- **Member**

  - Fields: `ID` (int), `Name` (string), `BorrowedBooks` (\[]Book)

### Service Interface (`LibraryManager`)

Defines system operations for:

- Book management (Add, Remove)
- Member registration (AddMember)
- Loan processing (BorrowBook, ReturnBook)
- Queries (ListAvailableBooks, ListBorrowedBooks)

---

## Controllers, Services, Models

### Controllers (CLI)

- **StartCLI** presents a console menu to the user.
- Input helpers:

  - `promptInt`: Prompts and validates numeric input.
  - `promptString`: Prompts and ensures non-empty input.

- Supported operations:

  1. Add book
  2. Remove book
  3. Add member
  4. Borrow book
  5. Return book
  6. List available books
  7. List books borrowed by a member
  8. Exit

---

### Services — Business Logic

Implemented in `Library` struct, utilizing:

```go
type Library struct {
  books   map[int]models.Book
  members map[int]*models.Member
}
```

Key methods:

- `AddBook` / `RemoveBook` — manage book catalog.
- `AddMember` — registers a member; prevents duplicates.
- `BorrowBook` — checks existence, status, member; updates status and logs borrowing.
- `ReturnBook` — validates return; updates status and removes from member’s list.
- `ListAvailableBooks` — filters books by availability.
- `ListBorrowedBooks` — fetches a member’s borrowed books.

---

### Models

#### `models/book.go`

```go
type Book struct {
  ID     int
  Title  string
  Author string
  Status string
}
```

#### `models/member.go`

```go
type Member struct {
  ID            int
  Name          string
  BorrowedBooks []Book
}
```

---

## Data Flow & User Interactions

1. **Launch**

   - `main.go` invokes `NewLibrary()` and `StartCLI(...)`.

2. **User Interaction**

   - CLI presents menu; users enter choices and data.

3. **Controller Actions**

   - Validates input, constructs data, calls corresponding service methods, and prints outcomes.

4. **Service Processing**

   - Updates in-memory data (`books`, `members`), returns results or errors.

---

## Error Handling & Validation

### Controllers

- Use `promptInt` and `promptString` to ensure valid input (non-empty strings, positive integers).
- Loop until valid input is entered or EOF.

### Services

- Return descriptive errors for:

  - Nonexistent book or member ID
  - Duplicate member registration
  - Attempted borrow of already borrowed book
  - Return of book not borrowed by the member

---

## Extensibility & Improvements

Possible future enhancements include:

- **Data persistence** (database or file storage)
- **Concurrency support** with thread safety
- **Web or API frontend**
- **Search & Filtering**, e.g., by title or author
- **Loan policies** like due dates, borrowing limits, fines
- **Member management** including updates and deletion

---

## Usage Instructions

1. **Run the application**

   ```bash
   go run main.go
   ```

2. **Choose from the menu:**

   - `1`: Add book (enter ID, title, author)
   - `2`: Remove book (enter ID)
   - `3`: Add member (enter ID, name)
   - `4`: Borrow book (enter book ID, member ID)
   - `5`: Return book (enter book ID, member ID)
   - `6`: List available books
   - `7`: List books borrowed by a member
   - `0`: Exit the program

---

## Example Scenarios

### Adding a Book

```
Enter book ID: 100
Enter title: The Hobbit
Enter author: J.R.R. Tolkien
Book added.
```

### Registering a Member

```
Enter member ID: 1
Enter name: Alice
Member added.
```

### Borrow & View Operations

```
Enter book ID to borrow: 100
Enter your member ID: 1
Book borrowed.

List available books:
No available books.

Books borrowed by member with ID 1:
#100: The Hobbit by J.R.R. Tolkien
```
