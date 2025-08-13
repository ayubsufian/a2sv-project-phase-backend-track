package repository

import (
	"context"
	"errors"
	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoUserRepository is a MongoDB-backed implementation of the UserRepository interface.
type mongoUserRepository struct {
	collection *mongo.Collection
}

// Add this compile-time check. It will fail to compile if method signatures don't match.
var _ usecase.IUserRepository = (*mongoUserRepository)(nil)

// NewMongoUserRepository initializes and returns a new user repository.
func NewMongoUserRepository(db *mongo.Database) usecase.IUserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

// Create inserts a new user document with a generated ObjectID.
func (r *mongoUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	oid := primitive.NewObjectID()
	doc := bson.D{
		{Key: "_id", Value: oid},
		{Key: "username", Value: u.Username},
		{Key: "password", Value: u.Password},
		{Key: "role", Value: u.Role},
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, usecase.ErrUserAlreadyExists
		}
		return domain.User{}, err
	}
	u.ID = oid.Hex()
	return u, nil
}

// FindByUsername looks up a user document by username.
func (r *mongoUserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var rec struct {
		ID       primitive.ObjectID `bson:"_id"`
		Username string             `bson:"username"`
		Password string             `bson:"password"`
		Role     string             `bson:"role"`
	}
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&rec)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, usecase.ErrNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:       rec.ID.Hex(),
		Username: rec.Username,
		Password: rec.Password,
		Role:     rec.Role,
	}, nil
}
