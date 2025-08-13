package usecase

import "errors"

var (
	// ErrUserAlreadyExists is returned from the Register use case when a user with the given username already exists in the system.
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrTaskAlreadyExists = errors.New("task already exists")

	ErrInvalidID = errors.New("invalid ID format")

	// ErrInvalidCredentials is returned from the Login use case when the provided password does not match the stored hash for the user.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrNotFound is a generic error returned when a requested resource (like a user or a task) cannot be found.
	ErrNotFound = errors.New("resource not found")
)
