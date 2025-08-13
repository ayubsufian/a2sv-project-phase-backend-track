package service

import (
	"testing"

	"task_manager_test/internal/usecase"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// PasswordServiceTestSuite defines the test suite for the bcryptHasher.
type PasswordServiceTestSuite struct {
	suite.Suite
	hasher usecase.IPasswordService
}

// SetupTest is run before each test in the suite and initializes the hasher.
func (s *PasswordServiceTestSuite) SetupTest() {
	s.hasher = NewPasswordHasher()
}

// TestPasswordService runs the entire test suite.
func TestPasswordService(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}

// TestHashSuccess tests the successful hashing of a password.
func (s *PasswordServiceTestSuite) TestHashSuccess() {
	password := "mySecretPassword"
	hashedPassword, err := s.hasher.Hash(password)

	// Assert that no error occurred during hashing
	s.Assert().NoError(err, "Hashing should not produce an error for a valid password")

	// Assert that the generated hash is not empty
	s.Assert().NotEmpty(hashedPassword, "The hashed password should not be empty")

	// Verify that the generated hash is a valid bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	s.Assert().NoError(err, "The generated hash should be verifiable against the original password")
}

// TestHashEmptyPassword tests the behavior of the Hash function with an empty password.
func (s *PasswordServiceTestSuite) TestHashEmptyPassword() {
	password := ""
	hashedPassword, err := s.hasher.Hash(password)

	// Assert that no error occurred during hashing
	s.Assert().NoError(err, "Hashing an empty password should not produce an error")

	// Assert that the generated hash is not empty
	s.Assert().NotEmpty(hashedPassword, "The hashed password for an empty string should not be empty")

	// Verify that the generated hash is a valid bcrypt hash for an empty string
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	s.Assert().NoError(err, "The generated hash should be verifiable against an empty password")
}

// TestCompareSuccess tests the successful comparison of a correct password and hash.
func (s *PasswordServiceTestSuite) TestCompareSuccess() {
	password := "correct-password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	s.Require().NoError(err, "Setup: Failed to generate hash for testing comparison")

	// Assert that the comparison returns true for the correct password
	match := s.hasher.Compare(string(hashedPassword), password)
	s.Assert().True(match, "Comparison should return true for a correct password")
}

// TestCompareFailure tests the failed comparison of an incorrect password and hash.
func (s *PasswordServiceTestSuite) TestCompareFailure() {
	password := "correct-password"
	incorrectPassword := "wrong-password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	s.Require().NoError(err, "Setup: Failed to generate hash for testing comparison")

	// Assert that the comparison returns false for an incorrect password
	match := s.hasher.Compare(string(hashedPassword), incorrectPassword)
	s.Assert().False(match, "Comparison should return false for an incorrect password")
}

// TestCompareInvalidHash tests the comparison with an invalid hash format.
func (s *PasswordServiceTestSuite) TestCompareInvalidHash() {
	invalidHash := "this-is-not-a-valid-bcrypt-hash"
	password := "any-password"

	// Assert that the comparison returns false for an invalid hash
	match := s.hasher.Compare(invalidHash, password)
	s.Assert().False(match, "Comparison should return false for an invalid hash format")
}
