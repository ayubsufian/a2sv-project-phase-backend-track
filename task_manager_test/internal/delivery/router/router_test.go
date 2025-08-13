package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"task_manager_test/internal/delivery/controller"
	"task_manager_test/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RouterTestSuite defines the test suite for the main application router.
type RouterTestSuite struct {
	suite.Suite
	router       *gin.Engine
	mockUserCont *controller.UserController
	mockTaskCont *controller.TaskController
	mockJwtSvc   *mocks.IJWTService
}

// getHandlerName retrieves the full function name for a given handler.
func getHandlerName(handler gin.HandlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}

// SetupTest is run before each test, setting up a fresh router configuration.
func (s *RouterTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUserCont = &controller.UserController{}
	s.mockTaskCont = &controller.TaskController{}
	s.mockJwtSvc = new(mocks.IJWTService)

	cfg := &RouterConfig{
		UserCont: s.mockUserCont,
		TaskCont: s.mockTaskCont,
		JwtSvc:   s.mockJwtSvc,
	}
	s.router = SetupRouter(cfg)
}

// TestRouter runs the entire test suite.
func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

// TestRouteRegistration verifies that all expected routes are registered correctly.
func (s *RouterTestSuite) TestRouteRegistration() {
	expectedRoutes := map[string]string{
		"POST:/register":           getHandlerName(s.mockUserCont.Register),
		"POST:/login":              getHandlerName(s.mockUserCont.Login),
		"GET:/api/tasks":           getHandlerName(s.mockTaskCont.GetTasks),
		"POST:/api/tasks":          getHandlerName(s.mockTaskCont.CreateTask),
		"GET:/api/tasks/:id":       getHandlerName(s.mockTaskCont.GetTask),
		"PUT:/api/tasks/:id":       getHandlerName(s.mockTaskCont.UpdateTask),
		"DELETE:/api/tasks/:id":    getHandlerName(s.mockTaskCont.DeleteTask),
		"GET:/api/admin/dashboard": getHandlerName(s.mockTaskCont.AdminDashboard),
	}

	registeredRoutes := s.router.Routes()
	actualRoutes := make(map[string]string)
	for _, route := range registeredRoutes {
		key := fmt.Sprintf("%s:%s", route.Method, route.Path)
		actualRoutes[key] = route.Handler
	}

	for expectedKey, expectedHandler := range expectedRoutes {
		actualHandler, ok := actualRoutes[expectedKey]
		assert.True(s.T(), ok, "Expected route %s to be registered", expectedKey)
		assert.True(s.T(), strings.HasSuffix(actualHandler, expectedHandler), "Route %s is registered with wrong handler. Expected %s, got %s", expectedKey, expectedHandler, actualHandler)
	}
}

// TestPublicRoutes checks that public routes are accessible without any authentication.
func (s *RouterTestSuite) TestPublicRoutes() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", nil)
	s.router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusBadRequest, w.Code, "Public routes should not be blocked by auth middleware")
}

// TestAuthMiddlewareIsApplied verifies that routes under the /api group are protected.
func (s *RouterTestSuite) TestAuthMiddlewareIsApplied() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	s.router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code, "Routes under /api should be protected")
	assert.JSONEq(s.T(), `{"error": "missing token"}`, w.Body.String(), "Should return a missing token error")
}

// TestAdminOnlyMiddlewareIsApplied verifies that routes under the /api/admin group are protected.
func (s *RouterTestSuite) TestAdminOnlyMiddlewareIsApplied() {
	validUserToken := "a-valid-user-token"
	userClaims := jwt.MapClaims{"username": "testuser", "role": "user"}
	s.mockJwtSvc.On("ValidateToken", validUserToken).Return(userClaims, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+validUserToken)
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code, "Routes under /api/admin should be protected by AdminOnly middleware")
	assert.JSONEq(s.T(), `{"error": "admin access required"}`, w.Body.String(), "Should return an admin access required error")
	s.mockJwtSvc.AssertExpectations(s.T())
}
