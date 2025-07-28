package usecases

import (
	"context"
	"task_manager_clean/domain"
	"task_manager_clean/repositories"
)

// TaskUsecase defines application-level operations for managing domain.Task entities.
type TaskUsecase interface {
	List(ctx context.Context) ([]domain.Task, error)
	Get(ctx context.Context, id string) (domain.Task, error)
	Create(ctx context.Context, t domain.Task) (domain.Task, error)
	Update(ctx context.Context, t domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id string) error
}

// taskUsecase implements TaskUsecase, orchestrating domain logic via TaskRepository.
type taskUsecase struct {
	repo repositories.TaskRepository
}

// NewTaskUsecase constructs a new TaskUsecase, injecting the repository dependency.
func NewTaskUsecase(repo repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{repo}
}

// List retrieves all domain.Task records using the repository layer.
func (u *taskUsecase) List(ctx context.Context) ([]domain.Task, error) {
	return u.repo.GetAll(ctx)
}

// Get fetches a task by its ID. Delegates error handling (e.g. invalid ID, missing record) to the repository.
func (u *taskUsecase) Get(ctx context.Context, id string) (domain.Task, error) {
	return u.repo.GetByID(ctx, id)
}

// Create builds and persists a new domain.Task entity.
func (u *taskUsecase) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	return u.repo.Create(ctx, t)
}

// Update modifies an existing domain.Task identified by its ID.
func (u *taskUsecase) Update(ctx context.Context, t domain.Task) (domain.Task, error) {
	return u.repo.Update(ctx, t)
}

// Delete removes a task by its ID. Business errors like 'not found' are surfaced from repository.
func (u *taskUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
