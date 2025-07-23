package data

import (
	"context"
	"errors"
	"task_manager_auth/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// userCollection is the MongoDB collection for storing user records.
var userCollection *mongo.Collection

// InitUserCollection initializes the userCollection variable by connecting to the "users" collection in the "taskdb" database.
func InitUserCollection() {
	userCollection = client.Database("taskdb").Collection("users")
}

// Register creates a new user in the users collection
func Register(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, _ := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
	if count > 0 {
		return errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	_, err = userCollection.InsertOne(ctx, user)
	return err
}

// Login verifies a user's credentials
func Login(username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}
	return user, nil
}
