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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// TaskControllerTestSuite defines the test suite for the TaskController.
type TaskControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mockUsecase    *mocks.TaskUsecase
	taskController *TaskController
	sampleTask     domain.Task
	sampleTime     time.Time
}

// SetupTest runs before each test in the suite.
func (s *TaskControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockUsecase = new(mocks.TaskUsecase)
	s.taskController = NewTaskController(s.mockUsecase)
	s.router = gin.Default()
	taskRoutes := s.router.Group("/tasks")
	{
		taskRoutes.GET("", s.taskController.GetTasks)
		taskRoutes.POST("", s.taskController.CreateTask)
		taskRoutes.GET("/:id", s.taskController.GetTask)
		taskRoutes.PUT("/:id", s.taskController.UpdateTask)
		taskRoutes.DELETE("/:id", s.taskController.DeleteTask)
	}
	s.router.GET("/admin/dashboard", s.taskController.AdminDashboard)

	s.sampleTime, _ = time.Parse(time.RFC3339, "2025-01-01T15:04:05Z")
	s.sampleTask = domain.Task{
		ID:          "task-123",
		Title:       "Sample Task",
		Description: "A description for the sample task.",
		DueDate:     s.sampleTime,
		Status:      "Pending",
	}
}

// TestTaskController runs the entire test suite.
func TestTaskController(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}

// --- GetTasks ---//
func (s *TaskControllerTestSuite) TestGetTasks_Success() {
	s.mockUsecase.On("List", mock.Anything).Return([]domain.Task{s.sampleTask}, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal([]TaskResponse{mapToTaskResponse(s.sampleTask)})
	s.JSONEq(string(expectedBody), w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestGetTasks_Error() {
	s.mockUsecase.On("List", mock.Anything).Return(nil, errors.New("database error")).Once()
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.JSONEq(`{"error": "could not retrieve tasks"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// --- GetTask ---//
func (s *TaskControllerTestSuite) TestGetTask_Success() {
	s.mockUsecase.On("Get", mock.Anything, "task-123").Return(s.sampleTask, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/task-123", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal(mapToTaskResponse(s.sampleTask))
	s.JSONEq(string(expectedBody), w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestGetTask_NotFound() {
	s.mockUsecase.On("Get", mock.Anything, "non-existent-id").Return(domain.Task{}, usecase.ErrNotFound).Once()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/non-existent-id", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.JSONEq(`{"error": "task not found"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestGetTask_InvalidID() {
	// Arrange: Mock the use case to return an invalid ID error.
	s.mockUsecase.On("Get", mock.Anything, "invalid-id-format").Return(domain.Task{}, usecase.ErrInvalidID).Once()

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/tasks/invalid-id-format", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.JSONEq(`{"error": "invalid task ID format"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// --- CreateTask ---//
func (s *TaskControllerTestSuite) TestCreateTask_Success() {
	taskToCreate := domain.Task{Title: s.sampleTask.Title, Description: s.sampleTask.Description, DueDate: s.sampleTask.DueDate, Status: s.sampleTask.Status}
	s.mockUsecase.On("Create", mock.Anything, taskToCreate).Return(s.sampleTask, nil).Once()
	body, _ := json.Marshal(gin.H{"title": taskToCreate.Title, "description": taskToCreate.Description, "duedate": taskToCreate.DueDate, "status": taskToCreate.Status})
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
	expectedBody, _ := json.Marshal(mapToTaskResponse(s.sampleTask))
	s.JSONEq(string(expectedBody), w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestCreateTask_BadRequestBinding() {
	body, _ := json.Marshal(gin.H{"description": "only a description"})
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockUsecase.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *TaskControllerTestSuite) TestCreateTask_Conflict() {
	// Arrange: Mock the use case to return an already exists error.
	taskToCreate := domain.Task{Title: s.sampleTask.Title, Description: s.sampleTask.Description, DueDate: s.sampleTask.DueDate, Status: s.sampleTask.Status}
	s.mockUsecase.On("Create", mock.Anything, taskToCreate).Return(domain.Task{}, usecase.ErrTaskAlreadyExists).Once()

	// Act
	body, _ := json.Marshal(gin.H{"title": taskToCreate.Title, "description": taskToCreate.Description, "duedate": taskToCreate.DueDate, "status": taskToCreate.Status})
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusConflict, w.Code)
	s.JSONEq(`{"error": "a task with these details already exists"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestCreateTask_UsecaseValidationError() {
	// Arrange: Mock the use case to return a generic validation error.
	taskToCreate := domain.Task{Title: "title", Description: "", DueDate: s.sampleTime, Status: "invalid"}
	validationError := errors.New("status must be one of 'pending', 'in-progress', or 'done'")
	s.mockUsecase.On("Create", mock.Anything, taskToCreate).Return(domain.Task{}, validationError).Once()

	// Act
	body, _ := json.Marshal(gin.H{"title": taskToCreate.Title, "description": taskToCreate.Description, "duedate": taskToCreate.DueDate, "status": taskToCreate.Status})
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert: This hits the default case which returns the raw error.
	s.Equal(http.StatusBadRequest, w.Code)
	s.JSONEq(`{"error": "status must be one of 'pending', 'in-progress', or 'done'"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// --- UpdateTask ---//
func (s *TaskControllerTestSuite) TestUpdateTask_Success() {
	s.mockUsecase.On("Update", mock.Anything, s.sampleTask).Return(s.sampleTask, nil).Once()
	body, _ := json.Marshal(gin.H{"title": s.sampleTask.Title, "description": s.sampleTask.Description, "duedate": s.sampleTask.DueDate, "status": s.sampleTask.Status})
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+s.sampleTask.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal(mapToTaskResponse(s.sampleTask))
	s.JSONEq(string(expectedBody), w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestUpdateTask_BadRequestBinding() {
	// Arrange: Send a request with a missing required field.
	body, _ := json.Marshal(gin.H{"description": "only a description"})
	req, _ := http.NewRequest(http.MethodPut, "/tasks/task-123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.mockUsecase.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *TaskControllerTestSuite) TestUpdateTask_NotFound() {
	taskToUpdate := s.sampleTask
	s.mockUsecase.On("Update", mock.Anything, taskToUpdate).Return(domain.Task{}, usecase.ErrNotFound).Once()
	body, _ := json.Marshal(gin.H{"title": taskToUpdate.Title, "description": taskToUpdate.Description, "duedate": taskToUpdate.DueDate, "status": taskToUpdate.Status})
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+taskToUpdate.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.JSONEq(`{"error": "task not found"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestUpdateTask_InvalidID() {
	// Arrange: Mock the use case to return an invalid ID error.
	taskToUpdate := s.sampleTask
	s.mockUsecase.On("Update", mock.Anything, taskToUpdate).Return(domain.Task{}, usecase.ErrInvalidID).Once()
	body, _ := json.Marshal(gin.H{"title": taskToUpdate.Title, "description": taskToUpdate.Description, "duedate": taskToUpdate.DueDate, "status": taskToUpdate.Status})
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+taskToUpdate.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.JSONEq(`{"error": "invalid task ID format"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// --- DeleteTask ---//
func (s *TaskControllerTestSuite) TestDeleteTask_Success() {
	s.mockUsecase.On("Delete", mock.Anything, "task-123").Return(nil).Once()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/task-123", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNoContent, w.Code)
	s.Empty(w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestDeleteTask_NotFound() {
	s.mockUsecase.On("Delete", mock.Anything, "non-existent-id").Return(usecase.ErrNotFound).Once()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/non-existent-id", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.JSONEq(`{"error": "task not found"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestDeleteTask_InvalidID() {
	// Arrange: Mock the use case to return an invalid ID error.
	s.mockUsecase.On("Delete", mock.Anything, "invalid-id-format").Return(usecase.ErrInvalidID).Once()

	// Act
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/invalid-id-format", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.JSONEq(`{"error": "invalid task ID format"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerTestSuite) TestDeleteTask_InternalError() {
	s.mockUsecase.On("Delete", mock.Anything, "task-123").Return(errors.New("some internal error")).Once()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/task-123", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.JSONEq(`{"error": "could not delete task"}`, w.Body.String())
	s.mockUsecase.AssertExpectations(s.T())
}

// --- AdminDashboard ---//
func (s *TaskControllerTestSuite) TestAdminDashboard() {
	req, _ := http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.JSONEq(`{"message": "Welcome Admin"}`, w.Body.String())
}
