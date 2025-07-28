package usecases

import (
	"context"
	"errors"
	"task_manager_clean/domain"
	"task_manager_clean/infrastructure"
	"task_manager_clean/repositories"
)

// UserUsecase defines the business logic operations related to user management.
type UserUsecase interface {
	Register(ctx context.Context, u domain.User) error
	Login(ctx context.Context, username, password string) (string, error)
}

// userUsecase is the concrete implementation of UserUsecase.
type userUsecase struct {
	repo       repositories.UserRepository
	pwdService infrastructure.PasswordHasher
	jwtService infrastructure.JWTService
}

// NewUserUsecase creates a new instance of userUsecase with dependencies injected.
func NewUserUsecase(repo repositories.UserRepository, pwd infrastructure.PasswordHasher, jwtSvc infrastructure.JWTService) UserUsecase {
	return &userUsecase{repo, pwd, jwtSvc}
}

// Register registers a new user by hashing their password and saving them in the repository.
func (u *userUsecase) Register(ctx context.Context, user domain.User) error {
	hashed, err := u.pwdService.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	if user.Role == "" {
		user.Role = "user"
	}
	_, err = u.repo.Create(ctx, user)
	return err
}

// Login validates user credentials and generates a JWT token if successful.
func (u *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	usr, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if !u.pwdService.Compare(usr.Password, password) {
		return "", errors.New("invalid username or password")
	}
	return u.jwtService.GenerateToken(usr.Username, usr.Role)
}
