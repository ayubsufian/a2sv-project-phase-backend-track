package repositories

import (
	"context"
	"errors"
	"task_manager_clean/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskRepository defines CRUD operations for domain.Task.
type TaskRepository interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, error)
	Create(ctx context.Context, t domain.Task) (domain.Task, error)
	Update(ctx context.Context, t domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id string) error
}

// mongoTaskRepo is a MongoDB-based implementation of TaskRepository.
type mongoTaskRepo struct {
	col *mongo.Collection
}

// NewMongoTaskRepository initializes a TaskRepository for a given Mongo collection.
func NewMongoTaskRepository(col *mongo.Collection) TaskRepository {
	return &mongoTaskRepo{col}
}

// GetAll retrieves all task documents from MongoDB and maps them to domain.Task.
func (r *mongoTaskRepo) GetAll(ctx context.Context) ([]domain.Task, error) {
	cur, err := r.col.Find(ctx, bson.M{})
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
func (r *mongoTaskRepo) GetByID(ctx context.Context, id string) (domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, errors.New("invalid task id format")
	}
	var rec struct {
		ID          primitive.ObjectID `bson:"_id"`
		Title       string             `bson:"title"`
		Description string             `bson:"description"`
		DueDate     time.Time          `bson:"duedate"`
		Status      string             `bson:"status"`
	}
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&rec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("task not found")
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
func (r *mongoTaskRepo) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	oid := primitive.NewObjectID()
	_, err := r.col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: oid},
		{Key: "title", Value: t.Title},
		{Key: "description", Value: t.Description},
		{Key: "duedate", Value: t.DueDate},
		{Key: "status", Value: t.Status},
	})
	if mongo.IsDuplicateKeyError(err) {
		return domain.Task{}, errors.New("task already exists")
	}
	t.ID = oid.Hex()
	return t, err
}

// Update replaces an existing task document with the data in domain.Task.
func (r *mongoTaskRepo) Update(ctx context.Context, t domain.Task) (domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return domain.Task{}, errors.New("invalid task id format")
	}
	res, err := r.col.ReplaceOne(ctx, bson.D{{Key: "_id", Value: oid}}, bson.D{
		{Key: "title", Value: t.Title},
		{Key: "description", Value: t.Description},
		{Key: "duedate", Value: t.DueDate},
		{Key: "status", Value: t.Status},
	})
	if err != nil {
		return domain.Task{}, err
	}
	if res.MatchedCount == 0 {
		return domain.Task{}, errors.New("task not found")
	}
	return t, nil
}

// Delete removes a task document by its ID.
func (r *mongoTaskRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task id format")
	}
	res, err := r.col.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
