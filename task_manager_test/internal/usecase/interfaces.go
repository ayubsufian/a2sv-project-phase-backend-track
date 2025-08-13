package usecase

import (
	"context"
	"task_manager_test/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// IUserRepository defines domain-centric user methods for creating and finding users.
type IUserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

// ITaskRepository defines CRUD operations for domain.Task.
type ITaskRepository interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, error)
	Create(ctx context.Context, t domain.Task) (domain.Task, error)
	Update(ctx context.Context, t domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id string) error
}

// IJWTService defines methods for generating and validating JWT tokens.
type IJWTService interface {
	GenerateToken(username, role string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
}

// IPasswordService defines methods for hashing and verifying passwords.
type IPasswordService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) bool
}
