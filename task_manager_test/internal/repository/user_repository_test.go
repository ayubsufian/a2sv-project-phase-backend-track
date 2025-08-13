package repository

import (
	"context"
	"os"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepositoryTestSuite defines the integration test suite for the user repository.
type UserRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	repository usecase.IUserRepository
}

// SetupSuite runs once before the entire suite starts. It's responsible for setting up the database connection using a dedicated test URI.
func (s *UserRepositoryTestSuite) SetupSuite() {
	// Load .env file, which should contain the test database URI.
	if err := godotenv.Load("../../.env"); err != nil {
		s.T().Log("No .env file found, proceeding with environment variables")
	}

	uri := os.Getenv("MONGODB_URI_TEST")
	if uri == "" {
		// Skip the suite if the test database is not configured to avoid panics.
		s.T().Skip("MONGODB_URI_TEST environment variable not set, skipping integration tests")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	assert.NoError(s.T(), err, "Failed to connect to MongoDB")

	s.client = client
	s.db = client.Database("userdb_test") // Use a dedicated test database for users
	s.collection = s.db.Collection("users_test")
}

// TearDownSuite runs once after all tests in the suite have finished.
func (s *UserRepositoryTestSuite) TearDownSuite() {
	if s.client != nil {
		// Clean up the entire test database after the suite is done
		err := s.db.Drop(context.Background())
		assert.NoError(s.T(), err, "Failed to drop test database")

		err = s.client.Disconnect(context.Background())
		assert.NoError(s.T(), err, "Failed to disconnect from MongoDB")
	}
}

// SetupTest runs before each individual test. It instantiates the repository and ensures the database schema (indexes) is correctly configured.
func (s *UserRepositoryTestSuite) SetupTest() {
	// Create a new repository instance for each test to ensure isolation.
	s.repository = NewMongoUserRepository(s.db)
	// We manually set the collection to our specific test collection.
	(s.repository.(*mongoUserRepository)).collection = s.collection

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := s.collection.Indexes().CreateOne(context.Background(), indexModel)
	assert.NoError(s.T(), err, "SetupTest: failed to create unique index on username")
}

// TearDownTest runs after each individual test. It is CRITICAL for ensuring test isolation by cleaning up any data created during the test.
func (s *UserRepositoryTestSuite) TearDownTest() {
	// Drop the collection to ensure a clean state for the next test.
	err := s.collection.Drop(context.Background())
	assert.NoError(s.T(), err, "Failed to drop test collection")
}

// TestUserRepositoryTestSuite is the entry point for the Go test runner.
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

// --- Test Cases ---

// TestCreateAndFindByUsername_Success tests the complete lifecycle of creating a user and then successfully retrieving them by their username.
func (s *UserRepositoryTestSuite) TestCreateAndFindByUsername_Success() {
	// ARRANGE
	ctx := context.Background()
	userToCreate := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:     "member",
	}

	// ACT - Create
	createdUser, err := s.repository.Create(ctx, userToCreate)

	// ASSERT - Create
	assert.NoError(s.T(), err, "Create should not return an error")
	assert.NotEmpty(s.T(), createdUser.ID, "Created user should have a non-empty ID")
	assert.True(s.T(), primitive.IsValidObjectID(createdUser.ID), "Created user ID should be a valid ObjectID")
	assert.Equal(s.T(), userToCreate.Username, createdUser.Username)

	// ACT - FindByUsername
	fetchedUser, err := s.repository.FindByUsername(ctx, "testuser")

	// ASSERT - FindByUsername
	assert.NoError(s.T(), err, "FindByUsername should not return an error for a valid user")
	assert.Equal(s.T(), createdUser, fetchedUser, "Fetched user should match the created user")
}

// TestCreate_Fails_When_UserAlreadyExists verifies that the repository correctly handles attempts to create a user with a username that is already taken.
func (s *UserRepositoryTestSuite) TestCreate_Fails_When_UserAlreadyExists() {
	// ARRANGE
	ctx := context.Background()
	// First, create an initial user.
	_, err := s.repository.Create(ctx, domain.User{Username: "existinguser", Password: "p1", Role: "r1"})
	assert.NoError(s.T(), err, "Setup: failed to create initial user")

	// Prepare a new user with the same username
	duplicateUser := domain.User{
		Username: "existinguser",
		Password: "p2",
		Role:     "r2",
	}

	// ACT: This call will now fail because the unique index exists.
	_, err = s.repository.Create(ctx, duplicateUser)

	// ASSERT
	assert.Error(s.T(), err, "Create should return an error for a duplicate username")
	assert.ErrorIs(s.T(), err, usecase.ErrUserAlreadyExists, "The error should be usecase.ErrUserAlreadyExists")
}

// TestFindByUsername_Fails_When_NotFound ensures the repository returns the correct error when searching for a user that does not exist.
func (s *UserRepositoryTestSuite) TestFindByUsername_Fails_When_NotFound() {
	// ARRANGE
	ctx := context.Background()
	nonExistentUsername := "ghost"

	// ACT
	_, err := s.repository.FindByUsername(ctx, nonExistentUsername)

	// ASSERT
	assert.Error(s.T(), err, "FindByUsername should return an error for a non-existent user")
	assert.ErrorIs(s.T(), err, usecase.ErrNotFound, "The error should be usecase.ErrNotFound")
}
