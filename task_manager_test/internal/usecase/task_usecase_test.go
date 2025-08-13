package usecase

import (
	"context"
	"errors"
	"task_manager_test/internal/domain"

	"task_manager_test/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TaskUsecaseTestSuite defines the test suite for the task use case.
type TaskUsecaseTestSuite struct {
	suite.Suite
	mockTaskRepo *mocks.ITaskRepository
	usecase      TaskUsecase
}

// SetupTest is a method from testify/suite. It runs before EACH test, ensuring a clean state by re-initializing the mock and the use case.
func (s *TaskUsecaseTestSuite) SetupTest() {
	// Create a new instance of the mock repository for each test.
	s.mockTaskRepo = mocks.NewITaskRepository(s.T())

	// Create a new instance of the use case, injecting our mock repository.
	s.usecase = NewTaskUsecase(s.mockTaskRepo)
}

// TestTaskUsecaseTestSuite is the Go test runner's entry point for this suite.
func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

// --- Test Cases for the Create Method ---

// TestCreate_Success tests the happy path for task creation.
func (s *TaskUsecaseTestSuite) TestCreate_Success() {
	// ARRANGE: Define inputs and set up mock expectations.
	ctx := context.Background()
	taskToCreate := domain.Task{Title: "A Valid Title", Description: "A description", DueDate: time.Now()}

	s.mockTaskRepo.On("Create", ctx, taskToCreate).Return(taskToCreate, nil)

	// ACT: Call the method we are testing.
	createdTask, err := s.usecase.Create(ctx, taskToCreate)

	// ASSERT: Verify the outcome.
	assert.NoError(s.T(), err, "Create should not return an error on success")
	assert.Equal(s.T(), taskToCreate, createdTask, "The created task should match the input task")
}

// TestCreate_Fails_When_TitleIsEmpty tests the specific business rule within the Create method.
func (s *TaskUsecaseTestSuite) TestCreate_Fails_When_TitleIsEmpty() {
	// ARRANGE
	ctx := context.Background()
	// Create a task with an invalid (empty) title.
	taskWithEmptyTitle := domain.Task{Title: " ", Description: "A description"}

	// ACT
	_, err := s.usecase.Create(ctx, taskWithEmptyTitle)

	// ASSERT
	assert.Error(s.T(), err, "Create should return an error for an empty title")
	assert.EqualError(s.T(), err, "task title cannot be empty", "The error message should be specific")

	s.mockTaskRepo.AssertNotCalled(s.T(), "Create")
}

// TestCreate_Fails_When_RepositoryFails tests how the use case handles an error from its dependency.
func (s *TaskUsecaseTestSuite) TestCreate_Fails_When_RepositoryFails() {
	// ARRANGE
	ctx := context.Background()
	taskToCreate := domain.Task{Title: "A Valid Title"}
	repoError := errors.New("database connection failed")

	// Configure the mock to return an error when Create is called.
	s.mockTaskRepo.On("Create", ctx, taskToCreate).Return(domain.Task{}, repoError)

	// ACT
	_, err := s.usecase.Create(ctx, taskToCreate)

	// ASSERT
	assert.Error(s.T(), err, "Create should propagate errors from the repository")
	assert.ErrorIs(s.T(), err, repoError, "The error should be the one returned by the repository")
}

// TestGet_Success tests the happy path for retrieving a single task.

func (s *TaskUsecaseTestSuite) TestGet_Success() {
	// ARRANGE
	ctx := context.Background()
	taskID := "task-123"
	expectedTask := domain.Task{ID: taskID, Title: "Test Task"}

	// Configure the mock to return the expected task when GetByID is called.
	s.mockTaskRepo.On("GetByID", ctx, taskID).Return(expectedTask, nil)

	// ACT
	actualTask, err := s.usecase.Get(ctx, taskID)

	// ASSERT
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedTask, actualTask)
}

// TestGet_Fails_When_NotFound tests the case where the repository returns ErrNotFound.
func (s *TaskUsecaseTestSuite) TestGet_Fails_When_NotFound() {
	// ARRANGE
	ctx := context.Background()
	taskID := "non-existent-id"

	// Configure the mock to return our application's standard ErrNotFound.
	s.mockTaskRepo.On("GetByID", ctx, taskID).Return(domain.Task{}, ErrNotFound)

	// ACT
	_, err := s.usecase.Get(ctx, taskID)

	// ASSERT
	assert.Error(s.T(), err)

	assert.ErrorIs(s.T(), err, ErrNotFound)
}

// TestDelete_Success tests the happy path for deleting a task.
func (s *TaskUsecaseTestSuite) TestDelete_Success() {
	// ARRANGE
	ctx := context.Background()
	taskID := "task-to-delete"

	// Configure the mock repository to return no error for the Delete operation.
	s.mockTaskRepo.On("Delete", ctx, taskID).Return(nil)

	// ACT
	err := s.usecase.Delete(ctx, taskID)

	// ASSERT
	assert.NoError(s.T(), err)
}

// TestDelete_Fails_When_RepositoryFails tests error propagation for the Delete operation.
func (s *TaskUsecaseTestSuite) TestDelete_Fails_When_RepositoryFails() {
	// ARRANGE
	ctx := context.Background()
	taskID := "task-to-delete"
	repoError := errors.New("permission denied")

	// Configure the mock repository to return an error.
	s.mockTaskRepo.On("Delete", ctx, taskID).Return(repoError)

	// ACT
	err := s.usecase.Delete(ctx, taskID)

	// ASSERT
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, repoError)
}
