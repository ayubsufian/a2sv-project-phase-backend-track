package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/mocks"
	"task_manager_test/internal/usecase"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// UserControllerTestSuite defines the test suite for the UserController.
type UserControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mockUsecase    *mocks.UserUsecase
	userController *UserController
}

// SetupTest runs before each test in the suite, ensuring a clean state.
func (s *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	// Initialize the mock and controller
	s.mockUsecase = new(mocks.UserUsecase)
	s.userController = NewUserController(s.mockUsecase)

	// Set up the router and define user routes
	s.router = gin.Default()
	userRoutes := s.router.Group("/users")
	{
		userRoutes.POST("/register", s.userController.Register)
		userRoutes.POST("/login", s.userController.Login)
	}
}

// TestUserController runs the entire test suite.
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

//--- Register Endpoint Tests ---//

// TestRegister_Success tests the successful registration of a new user.
func (s *UserControllerTestSuite) TestRegister_Success() {
	// Arrange
	userToRegister := domain.User{Username: "testuser", Password: "password123", Role: "user"}
	s.mockUsecase.On("Register", mock.Anything, userToRegister).Return(nil).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "testuser", "password": "password123", "role": "user"})
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.JSONEq(`{"message": "User Registered successfully"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestRegister_BadRequest tests a registration attempt with invalid JSON body.
func (s *UserControllerTestSuite) TestRegister_BadRequest() {
	// Arrange: Body is missing the required "password" field.
	body, _ := json.Marshal(gin.H{"username": "testuser", "role": "user"})
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	// Verify that the use case was never called because binding failed first.
	s.mockUsecase.AssertNotCalled(s.T(), "Register", mock.Anything, mock.Anything)
}

// TestRegister_Conflict tests a registration attempt where the user already exists.
func (s *UserControllerTestSuite) TestRegister_Conflict() {
	// Arrange
	userToRegister := domain.User{Username: "existinguser", Password: "password123", Role: "user"}
	s.mockUsecase.On("Register", mock.Anything, userToRegister).Return(usecase.ErrUserAlreadyExists).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "existinguser", "password": "password123", "role": "user"})
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusConflict, w.Code)
	s.JSONEq(`{"error": "a user with this username already exists"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestRegister_InternalError tests an unexpected server error during registration.
func (s *UserControllerTestSuite) TestRegister_InternalError() {
	// Arrange
	userToRegister := domain.User{Username: "testuser", Password: "password123", Role: "user"}
	// This hits the default case in the controller's switch statement.
	s.mockUsecase.On("Register", mock.Anything, userToRegister).Return(errors.New("database connection failed")).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "testuser", "password": "password123", "role": "user"})
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	s.JSONEq(`{"error": "could not register user"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

//--- Login Endpoint Tests ---//

// TestLogin_Success tests a successful user login.
func (s *UserControllerTestSuite) TestLogin_Success() {
	// Arrange
	expectedToken := "a-valid-jwt-token"
	s.mockUsecase.On("Login", mock.Anything, "testuser", "password123").Return(expectedToken, nil).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "testuser", "password": "password123"})
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.JSONEq(`{"token": "a-valid-jwt-token"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestLogin_BadRequest tests a login attempt with an invalid JSON body.
func (s *UserControllerTestSuite) TestLogin_BadRequest() {
	// Arrange: Body is missing the required "password" field.
	body, _ := json.Marshal(gin.H{"username": "testuser"})
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.mockUsecase.AssertNotCalled(s.T(), "Login", mock.Anything, mock.Anything, mock.Anything)
}

// TestLogin_Unauthorized_NotFound tests a login attempt for a user that does not exist.
func (s *UserControllerTestSuite) TestLogin_Unauthorized_NotFound() {
	// Arrange
	s.mockUsecase.On("Login", mock.Anything, "nonexistent", "password123").Return("", usecase.ErrNotFound).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "nonexistent", "password": "password123"})
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.JSONEq(`{"error": "invalid username or password"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestLogin_Unauthorized_InvalidCredentials tests a login attempt with an incorrect password.
func (s *UserControllerTestSuite) TestLogin_Unauthorized_InvalidCredentials() {
	// Arrange
	s.mockUsecase.On("Login", mock.Anything, "testuser", "wrongpassword").Return("", usecase.ErrInvalidCredentials).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "testuser", "password": "wrongpassword"})
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.JSONEq(`{"error": "invalid username or password"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestLogin_InternalError tests an unexpected server error during login.
func (s *UserControllerTestSuite) TestLogin_InternalError() {
	// Arrange
	// This hits the default case in the controller's switch statement.
	s.mockUsecase.On("Login", mock.Anything, "testuser", "password123").Return("", errors.New("token generation failed")).Once()

	// Act
	body, _ := json.Marshal(gin.H{"username": "testuser", "password": "password123"})
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	s.JSONEq(`{"error": "an internal server error occurred"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}
