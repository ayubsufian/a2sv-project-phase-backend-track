package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"task_manager_test/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AuthMiddlewareTestSuite defines the test suite for auth-related middleware.
type AuthMiddlewareTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mockJWTService *mocks.IJWTService
}

// SetupTest is run before each test in the suite.
func (s *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockJWTService = new(mocks.IJWTService)
	s.router = gin.New()
}

// TestAuthMiddleware runs the entire test suite.
func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

// TestAuthMiddleware_Success tests the happy path where a valid token is provided.
func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_Success() {
	// Arrange
	validToken := "valid.jwt.token"

	expectedClaims := jwt.MapClaims{"username": "testuser", "role": "admin"}
	s.mockJWTService.On("ValidateToken", validToken).Return(expectedClaims, nil).Once()

	// Apply middleware to a test route
	s.router.GET("/protected", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		s.Equal("testuser", username)
		s.Equal("admin", role)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.JSONEq(`{"status": "ok"}`, w.Body.String())
	s.mockJWTService.AssertExpectations(s.T())
}

// TestAuthMiddleware_NoAuthHeader tests the case where the Authorization header is missing.
func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_NoAuthHeader() {
	s.router.GET("/protected", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.JSONEq(`{"error": "missing token"}`, w.Body.String())
	s.mockJWTService.AssertNotCalled(s.T(), "ValidateToken", mock.Anything)
}

// TestAuthMiddleware_BadHeaderFormat tests for a malformed Authorization header.
func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_BadHeaderFormat() {
	s.router.GET("/protected", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic some-other-token-format")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.JSONEq(`{"error": "missing token"}`, w.Body.String())
	s.mockJWTService.AssertNotCalled(s.T(), "ValidateToken", mock.Anything)
}

// TestAuthMiddleware_InvalidToken tests for a token that is rejected by the JWT service.
func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	invalidToken := "invalid.or.expired.token"
	s.mockJWTService.On("ValidateToken", invalidToken).Return(nil, errors.New("token is expired")).Once()

	s.router.GET("/protected", AuthMiddleware(s.mockJWTService), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.JSONEq(`{"error": "invalid or expired token"}`, w.Body.String())
	s.mockJWTService.AssertExpectations(s.T())
}

//--- AdminOnly Middleware Tests ---//

// TestAdminOnly_Success tests when an admin user tries to access a restricted route.
func (s *AuthMiddlewareTestSuite) TestAdminOnly_Success() {
	s.router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	}, AdminOnly(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "welcome admin"})
	})
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.JSONEq(`{"status": "welcome admin"}`, w.Body.String())
}

// TestAdminOnly_Forbidden_NotAdmin tests when a non-admin user tries to access a restricted route.
func (s *AuthMiddlewareTestSuite) TestAdminOnly_Forbidden_NotAdmin() {
	s.router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	}, AdminOnly(), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusForbidden, w.Code)
	s.JSONEq(`{"error": "admin access required"}`, w.Body.String())
}

// TestAdminOnly_Forbidden_NoRole tests when the role is not set in the context at all.
func (s *AuthMiddlewareTestSuite) TestAdminOnly_Forbidden_NoRole() {
	s.router.GET("/admin", AdminOnly(), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusForbidden, w.Code)
	s.JSONEq(`{"error": "admin access required"}`, w.Body.String())
}

// TestAdminOnly_Forbidden_WrongRoleType tests when the role has an unexpected type.
func (s *AuthMiddlewareTestSuite) TestAdminOnly_Forbidden_WrongRoleType() {
	s.router.GET("/admin", func(c *gin.Context) {
		c.Set("role", 123)
		c.Next()
	}, AdminOnly(), func(c *gin.Context) {
		s.Fail("Next handler should not be called")
	})
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusForbidden, w.Code)
	s.JSONEq(`{"error": "admin access required"}`, w.Body.String())
}
