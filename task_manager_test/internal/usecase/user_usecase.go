package usecase

import (
	"context"
	"task_manager_test/internal/domain"
)

// UserUsecase defines the business logic operations related to user management.
type UserUsecase interface {
	Register(ctx context.Context, u domain.User) error
	Login(ctx context.Context, username, password string) (string, error)
}

// userUsecase is the concrete implementation of UserUsecase.
type userUsecase struct {
	repo       IUserRepository
	pwdService IPasswordService
	jwtService IJWTService
}

// NewUserUsecase creates a new instance of userUsecase with dependencies injected.
func NewUserUsecase(repo IUserRepository, pwd IPasswordService, jwtSvc IJWTService) UserUsecase {
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
		return "", ErrInvalidCredentials
	}
	return u.jwtService.GenerateToken(usr.Username, usr.Role)
}
