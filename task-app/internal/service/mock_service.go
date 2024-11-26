package service

import (
	"context"
	"task-app/internal/models"

	"github.com/google/uuid"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
	ListTasks(ctx context.Context, filters map[string]interface{}) ([]models.Task, error)
	DuplicateTask(ctx context.Context, id uuid.UUID) (*models.Task, error)
}
