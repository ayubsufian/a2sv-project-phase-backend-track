package repository

import (
	"context"
	"os"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskRepositoryTestSuite defines the integration test suite for the task repository.
type TaskRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	repository usecase.ITaskRepository
}

// SetupSuite runs once before the entire suite starts. It's responsible for setting up the database connection.
func (s *TaskRepositoryTestSuite) SetupSuite() {
	// Load .env file, which should contain the test database URI
	if err := godotenv.Load("../../.env"); err != nil {
		s.T().Log("No .env file found, proceeding with environment variables")
	}

	uri := os.Getenv("MONGODB_URI_TEST")
	if uri == "" {
		// Skip the suite if the test database is not configured.
		s.T().Skip("MONGODB_URI_TEST environment variable not set, skipping integration tests")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	assert.NoError(s.T(), err, "Failed to connect to MongoDB")

	s.client = client
	s.db = client.Database("taskdb_test") // Use a dedicated test database
	s.collection = s.db.Collection("tasks_test")
}

// TearDownSuite runs once after all tests in the suite have finished.
func (s *TaskRepositoryTestSuite) TearDownSuite() {
	if s.client != nil {
		err := s.client.Disconnect(context.Background())
		assert.NoError(s.T(), err, "Failed to disconnect from MongoDB")
	}
}

// SetupTest runs before each individual test. It instantiates the repository.
func (s *TaskRepositoryTestSuite) SetupTest() {
	s.repository = NewMongoTaskRepository(s.db)
	(s.repository.(*mongoTaskRepository)).collection = s.collection
}

// TearDownTest runs after each individual test. It's CRITICAL for ensuring test isolation by cleaning up any data created during the test.
func (s *TaskRepositoryTestSuite) TearDownTest() {
	// Drop the collection to ensure a clean state for the next test.
	err := s.collection.Drop(context.Background())
	assert.NoError(s.T(), err, "Failed to drop test collection")
}

// TestTaskRepositoryTestSuite is the entry point for the Go test runner.
func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}

// --- Test Cases ---

func (s *TaskRepositoryTestSuite) TestCreateAndGetByID_Success() {
	// ARRANGE
	ctx := context.Background()
	taskToCreate := domain.Task{
		Title:       "Integration Test Task",
		Description: "A task created during an integration test.",
		DueDate:     time.Now().Add(24 * time.Hour).UTC().Truncate(time.Millisecond),
		Status:      "To Do",
	}

	// ACT - Create
	createdTask, err := s.repository.Create(ctx, taskToCreate)

	// ASSERT - Create
	assert.NoError(s.T(), err, "Create should not return an error")
	assert.NotEmpty(s.T(), createdTask.ID, "Created task should have a non-empty ID")
	assert.True(s.T(), primitive.IsValidObjectID(createdTask.ID), "Created task ID should be a valid ObjectID")
	// Verify other fields were set correctly
	assert.Equal(s.T(), taskToCreate.Title, createdTask.Title)
	assert.Equal(s.T(), taskToCreate.Description, createdTask.Description)
	assert.Equal(s.T(), taskToCreate.Status, createdTask.Status)
	assert.WithinDuration(s.T(), taskToCreate.DueDate, createdTask.DueDate, time.Millisecond)

	// ACT - GetByID
	fetchedTask, err := s.repository.GetByID(ctx, createdTask.ID)

	// ASSERT - GetByID
	assert.NoError(s.T(), err, "GetByID should not return an error for a valid ID")
	assert.Equal(s.T(), createdTask, fetchedTask, "Fetched task should match the created task")
}

func (s *TaskRepositoryTestSuite) TestGetByID_Fails_When_NotFound() {
	// ARRANGE
	ctx := context.Background()
	nonExistentID := primitive.NewObjectID().Hex()

	// ACT
	_, err := s.repository.GetByID(ctx, nonExistentID)

	// ASSERT
	assert.Error(s.T(), err, "GetByID should return an error for a non-existent ID")
	assert.ErrorIs(s.T(), err, usecase.ErrNotFound, "The error should be usecase.ErrNotFound")
}

func (s *TaskRepositoryTestSuite) TestGetByID_Fails_When_InvalidIDFormat() {
	// ARRANGE
	ctx := context.Background()
	invalidID := "this-is-not-a-valid-object-id"

	// ACT
	_, err := s.repository.GetByID(ctx, invalidID)

	// ASSERT
	assert.Error(s.T(), err, "GetByID should return an error for an invalid ID format")
	assert.ErrorIs(s.T(), err, usecase.ErrInvalidID, "The error should be usecase.ErrInvalidID")
}

func (s *TaskRepositoryTestSuite) TestUpdate_Success() {
	// ARRANGE - Create a task first
	ctx := context.Background()
	initialTask, _ := s.repository.Create(ctx, domain.Task{Title: "Initial Title", Status: "To Do"})

	// ACT - Update the task
	updatedTaskData := domain.Task{
		ID:          initialTask.ID,
		Title:       "Updated Title",
		Description: "This description has been added.",
		Status:      "In Progress",
		DueDate:     initialTask.DueDate,
	}
	resultTask, err := s.repository.Update(ctx, updatedTaskData)

	// ASSERT
	assert.NoError(s.T(), err, "Update should not return an error")
	assert.Equal(s.T(), updatedTaskData, resultTask, "Result task should match the updated data")

	// Verify by fetching from the DB again
	fetchedTask, _ := s.repository.GetByID(ctx, initialTask.ID)
	assert.Equal(s.T(), "Updated Title", fetchedTask.Title)
	assert.Equal(s.T(), "In Progress", fetchedTask.Status)
}

func (s *TaskRepositoryTestSuite) TestDelete_Success() {
	// ARRANGE - Create a task to delete
	ctx := context.Background()
	taskToDelete, _ := s.repository.Create(ctx, domain.Task{Title: "Task to be Deleted"})

	// ACT - Delete the task
	err := s.repository.Delete(ctx, taskToDelete.ID)

	// ASSERT
	assert.NoError(s.T(), err, "Delete should not return an error for an existing task")

	// Verify by trying to fetch it again, which should now fail
	_, fetchErr := s.repository.GetByID(ctx, taskToDelete.ID)
	assert.ErrorIs(s.T(), fetchErr, usecase.ErrNotFound, "Fetching a deleted task should result in ErrNotFound")
}

func (s *TaskRepositoryTestSuite) TestDelete_Fails_When_NotFound() {
	// ARRANGE
	ctx := context.Background()
	nonExistentID := primitive.NewObjectID().Hex()

	// ACT
	err := s.repository.Delete(ctx, nonExistentID)

	// ASSERT
	assert.Error(s.T(), err, "Delete should return an error for a non-existent ID")
	assert.ErrorIs(s.T(), err, usecase.ErrNotFound, "The error should be usecase.ErrNotFound")
}
