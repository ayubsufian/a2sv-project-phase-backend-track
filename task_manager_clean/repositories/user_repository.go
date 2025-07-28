package repositories

import (
	"context"
	"errors"
	"task_manager_clean/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository defines domain-centric user methods for creating and finding users.
type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

// mongoUserRepo is a MongoDB-backed implementation of UserRepository.
type mongoUserRepo struct {
	col *mongo.Collection
}

// NewMongoUserRepository initializes and returns a new mongoUserRepo.
func NewMongoUserRepository(col *mongo.Collection) UserRepository {
	return &mongoUserRepo{col}
}

// Create inserts a new user document with a generated ObjectID.
func (r *mongoUserRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	oid := primitive.NewObjectID()
	doc := bson.D{
		{Key: "_id", Value: oid},
		{Key: "username", Value: u.Username},
		{Key: "password", Value: u.Password},
		{Key: "role", Value: u.Role},
	}

	_, err := r.col.InsertOne(ctx, doc)
	if mongo.IsDuplicateKeyError(err) {
		return domain.User{}, errors.New("username already exists")
	}
	if err != nil {
		return domain.User{}, err
	}
	u.ID = oid.Hex()
	return u, nil
}

// FindByUsername looks up a user document by username.
func (r *mongoUserRepo) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var rec struct {
		ID       primitive.ObjectID `bson:"_id"`
		Username string             `bson:"username"`
		Password string             `bson:"password"`
		Role     string             `bson:"role"`
	}
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&rec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.New("invalid username or password")
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
