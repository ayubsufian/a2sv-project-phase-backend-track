package repository

import (
	"context"
	"errors"

	"task_manager_test/internal/domain"
	"task_manager_test/internal/usecase"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoTaskRepository is the MongoDB-based implementation of the TaskRepository interface.
type mongoTaskRepository struct {
	collection *mongo.Collection
}

// Add a compile-time check to ensure this struct implements the correct interface.
var _ usecase.ITaskRepository = (*mongoTaskRepository)(nil)

// NewMongoTaskRepository is the constructor for the implementation.
func NewMongoTaskRepository(db *mongo.Database) usecase.ITaskRepository {
	return &mongoTaskRepository{
		collection: db.Collection("tasks"),
	}
}

// GetAll retrieves all task documents from MongoDB and maps them to domain.Task.
func (r *mongoTaskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []domain.Task
	for cur.Next(ctx) {
		var rec struct {
			ID          primitive.ObjectID `bson:"_id"`
			Title       string             `bson:"title"`
			Description string             `bson:"description"`
			DueDate     time.Time          `bson:"duedate"`
			Status      string             `bson:"status"`
		}
		if err := cur.Decode(&rec); err != nil {
			return nil, err
		}
		out = append(out, domain.Task{
			ID:          rec.ID.Hex(),
			Title:       rec.Title,
			Description: rec.Description,
			DueDate:     rec.DueDate,
			Status:      rec.Status,
		})
	}
	return out, nil
}

// GetByID fetches a task by its hexadecimal string ID.
func (r *mongoTaskRepository) GetByID(ctx context.Context, id string) (domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, usecase.ErrInvalidID
	}
	var rec struct {
		ID          primitive.ObjectID `bson:"_id"`
		Title       string             `bson:"title"`
		Description string             `bson:"description"`
		DueDate     time.Time          `bson:"duedate"`
		Status      string             `bson:"status"`
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&rec)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Task{}, usecase.ErrNotFound
		}
		return domain.Task{}, err
	}

	return domain.Task{
		ID:          rec.ID.Hex(),
		Title:       rec.Title,
		Description: rec.Description,
		DueDate:     rec.DueDate,
		Status:      rec.Status,
	}, nil
}

// Create inserts a new task document, generating a new unique ID.
func (r *mongoTaskRepository) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	oid := primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: oid},
		{Key: "title", Value: t.Title},
		{Key: "description", Value: t.Description},
		{Key: "duedate", Value: t.DueDate},
		{Key: "status", Value: t.Status},
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Task{}, usecase.ErrTaskAlreadyExists
		}
		return domain.Task{}, err
	}
	t.ID = oid.Hex()
	return t, err
}

// Update replaces an existing task document with the data in domain.Task.
func (r *mongoTaskRepository) Update(ctx context.Context, t domain.Task) (domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return domain.Task{}, usecase.ErrInvalidID
	}
	res, err := r.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: oid}}, bson.D{
		{Key: "title", Value: t.Title},
		{Key: "description", Value: t.Description},
		{Key: "duedate", Value: t.DueDate},
		{Key: "status", Value: t.Status},
	})
	if err != nil {
		return domain.Task{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Task{}, usecase.ErrNotFound
	}
	return t, nil
}

// Delete removes a task document by its ID.
func (r *mongoTaskRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return usecase.ErrInvalidID
	}
	res, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return usecase.ErrNotFound
	}
	return nil
}
