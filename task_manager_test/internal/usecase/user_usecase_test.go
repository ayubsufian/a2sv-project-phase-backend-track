package usecase

import (
	"context"
	"errors"
	"task_manager_test/internal/domain"

	"task_manager_test/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserUsecaseTestSuite defines the test suite for the user use case.
type UserUsecaseTestSuite struct {
	suite.Suite
	mockUserRepo *mocks.IUserRepository
	mockPwdSvc   *mocks.IPasswordService
	mockJwtSvc   *mocks.IJWTService
	usecase      UserUsecase
}

// SetupTest is a method from testify/suite. It runs before EACH test in the suite.
func (s *UserUsecaseTestSuite) SetupTest() {
	// Create new instances of our mocks for every single test.
	s.mockUserRepo = mocks.NewIUserRepository(s.T())
	s.mockPwdSvc = mocks.NewIPasswordService(s.T())
	s.mockJwtSvc = mocks.NewIJWTService(s.T())

	// Create a new instance of the use case we're testing, injecting our mock dependencies.
	s.usecase = NewUserUsecase(s.mockUserRepo, s.mockPwdSvc, s.mockJwtSvc)
}

// TestUserUsecaseTestSuite is the Go test runner's entry point for this suite.
func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

// --- Test Cases for the Register Method ---

// TestRegister_Success tests the "happy path" for user registration.
func (s *UserUsecaseTestSuite) TestRegister_Success() {
	// ARRANGE: Define inputs and set up mock expectations.
	ctx := context.Background()
	plainPassword := "plain-password"
	hashedPassword := "hashed-password"
	userToRegister := domain.User{Username: "newuser", Password: plainPassword, Role: "user"}

	s.mockPwdSvc.On("Hash", plainPassword).Return(hashedPassword, nil)

	expectedUserInRepo := domain.User{Username: "newuser", Password: hashedPassword, Role: "user"}
	s.mockUserRepo.On("Create", ctx, expectedUserInRepo).Return(expectedUserInRepo, nil)

	// ACT: Call the actual method we are testing.
	err := s.usecase.Register(ctx, userToRegister)

	// ASSERT: Verify the results.
	assert.NoError(s.T(), err, "Register should not return an error on success")
}

// TestRegister_Fails_When_UserAlreadyExists tests the case where the repository indicates that the username is already taken.
func (s *UserUsecaseTestSuite) TestRegister_Fails_When_UserAlreadyExists() {
	// ARRANGE
	ctx := context.Background()
	userToRegister := domain.User{Username: "existinguser", Password: "password", Role: "user"}

	s.mockPwdSvc.On("Hash", "password").Return("hashed-password", nil)

	expectedUserInRepo := domain.User{Username: "existinguser", Password: "hashed-password", Role: "user"}
	s.mockUserRepo.On("Create", ctx, expectedUserInRepo).Return(domain.User{}, ErrUserAlreadyExists)

	// ACT
	err := s.usecase.Register(ctx, userToRegister)

	// ASSERT
	assert.Error(s.T(), err, "Register should return an error when user exists")
	assert.ErrorIs(s.T(), err, ErrUserAlreadyExists, "The returned error should be ErrUserAlreadyExists")
}

// TestRegister_Fails_When_HashingFails tests an infrastructure failure scenario.
func (s *UserUsecaseTestSuite) TestRegister_Fails_When_HashingFails() {
	// ARRANGE
	ctx := context.Background()
	userToRegister := domain.User{Username: "anyuser", Password: "password"}
	hashingError := errors.New("bcrypt cost is too high")

	s.mockPwdSvc.On("Hash", "password").Return("", hashingError)

	// ACT
	err := s.usecase.Register(ctx, userToRegister)

	// ASSERT
	assert.Error(s.T(), err, "Register should return an error if hashing fails")
	assert.ErrorIs(s.T(), err, hashingError, "The error from the hashing service should be propagated")
	s.mockUserRepo.AssertNotCalled(s.T(), "Create")
}

// --- Test Cases for the Login Method ---

// TestLogin_Success tests the "happy path" for user login.
func (s *UserUsecaseTestSuite) TestLogin_Success() {
	// ARRANGE
	ctx := context.Background()
	username := "testuser"
	plainPassword := "correct-password"
	hashedPassword := "hashed-password"
	role := "admin"
	expectedToken := "a-valid-jwt-token"
	userFromRepo := domain.User{ID: "user-123", Username: username, Password: hashedPassword, Role: role}

	s.mockUserRepo.On("FindByUsername", ctx, username).Return(userFromRepo, nil)
	s.mockPwdSvc.On("Compare", hashedPassword, plainPassword).Return(true)
	s.mockJwtSvc.On("GenerateToken", username, role).Return(expectedToken, nil)

	// ACT
	token, err := s.usecase.Login(ctx, username, plainPassword)

	// ASSERT
	assert.NoError(s.T(), err, "Login should not return an error on success")
	assert.Equal(s.T(), expectedToken, token, "The returned token should match the expected token")
}

// TestLogin_Fails_When_UserNotFound tests the scenario where the username does not exist.
func (s *UserUsecaseTestSuite) TestLogin_Fails_When_UserNotFound() {
	// ARRANGE
	ctx := context.Background()
	username := "non-existent-user"

	s.mockUserRepo.On("FindByUsername", ctx, username).Return(domain.User{}, ErrNotFound)

	// ACT
	token, err := s.usecase.Login(ctx, username, "any-password")

	// ASSERT
	assert.Error(s.T(), err, "Login should return an error when user is not found")
	assert.ErrorIs(s.T(), err, ErrNotFound, "The error should be ErrNotFound")
	assert.Empty(s.T(), token, "No token should be returned on failure")
	// Assert that subsequent services were never called.
	s.mockPwdSvc.AssertNotCalled(s.T(), "Compare")
	s.mockJwtSvc.AssertNotCalled(s.T(), "GenerateToken")
}

// TestLogin_Fails_When_PasswordIsIncorrect tests when the password does not match.
func (s *UserUsecaseTestSuite) TestLogin_Fails_When_PasswordIsIncorrect() {
	// ARRANGE
	ctx := context.Background()
	username := "testuser"
	wrongPassword := "wrong-password"
	hashedPassword := "hashed-password"
	userFromRepo := domain.User{ID: "user-123", Username: username, Password: hashedPassword}

	s.mockUserRepo.On("FindByUsername", ctx, username).Return(userFromRepo, nil)
	s.mockPwdSvc.On("Compare", hashedPassword, wrongPassword).Return(false)

	// ACT
	token, err := s.usecase.Login(ctx, username, wrongPassword)

	// ASSERT
	assert.Error(s.T(), err, "Login should return an error for invalid credentials")
	assert.ErrorIs(s.T(), err, ErrInvalidCredentials, "The error should be ErrInvalidCredentials")
	assert.Empty(s.T(), token, "No token should be returned on failure")
	// Assert that the JWT service was never called.
	s.mockJwtSvc.AssertNotCalled(s.T(), "GenerateToken")
}
