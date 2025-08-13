package service

import (
	"task_manager_test/internal/usecase"

	"golang.org/x/crypto/bcrypt"
)

// bcryptHasher is an implementation of the PasswordService interface using bcrypt.
type bcryptHasher struct{}

// This compile-time check ensures that *bcryptHasher satisfies the IPasswordService interface.
var _ usecase.IPasswordService = (*bcryptHasher)(nil)

// NewPasswordHasher constructs a new instance of bcryptHasher.
func NewPasswordHasher() usecase.IPasswordService {
	return &bcryptHasher{}
}

// Hash generates a bcrypt hash from a plain-text password.
func (h *bcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// Compare verifies whether the plain-text password matches the bcrypt hash.
func (h *bcryptHasher) Compare(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
