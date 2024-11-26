package repository

import (
	"context"
	"github.com/google/uuid"
	"task-app/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}) ([]models.Task, error)
	Duplicate(ctx context.Context, id uuid.UUID) (*models.Task, error)
}
