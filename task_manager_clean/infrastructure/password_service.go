package infrastructure

import "golang.org/x/crypto/bcrypt"

// PasswordHasher defines methods for hashing and verifying passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) bool
}

// bcryptHasher is an implementation of PasswordHasher using bcrypt.
type bcryptHasher struct{}

// NewPasswordHasher constructs a new instance of bcryptHasher.
func NewPasswordHasher() PasswordHasher {
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
