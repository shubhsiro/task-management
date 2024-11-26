package service

import (
	"context"
	"log"
	"task-app/internal/models"
	"task-app/internal/repository"
	"task-app/pkg/utils"

	"github.com/google/uuid"
)

type TaskService struct {
	repo      repository.TaskRepository
	validator *utils.Validator
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo:      repo,
		validator: utils.NewValidator(),
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error) {
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Status:      req.Status,
	}

	log.Printf("Creating task: Title=%s, Category=%s, Status=%s", task.Title, task.Category, task.Status)

	if err := s.validator.ValidateTask(task); err != nil {
		log.Printf("Task validation failed: %v", err)

		return nil, err
	}

	err := s.repo.Create(ctx, task)
	if err != nil {
		log.Printf("Failed to create task: Title=%s, Category=%s, Error=%v", task.Title, task.Category, err)

		return nil, err
	}

	log.Printf("Task created successfully: ID=%s", task.ID)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	log.Printf("Retrieving task: ID=%s", id)

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Failed to retrieve task: ID=%s, Error=%v", id, err)

		return nil, err
	}

	log.Printf("Task retrieved successfully: ID=%s, Title=%s", task.ID, task.Title)
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	log.Printf("Updating task: ID=%s, Title=%s, Category=%s, Status=%s", task.ID, task.Title, task.Category, task.Status)

	if err := s.validator.ValidateTask(task); err != nil {
		log.Printf("Task validation failed: ID=%s, Error=%v", task.ID, err)
		return err
	}

	err := s.repo.Update(ctx, task)
	if err != nil {
		log.Printf("Failed to update task: ID=%s, Error=%v", task.ID, err)

		return err
	}

	log.Printf("Task updated successfully: ID=%s, Title=%s", task.ID, task.Title)

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting task: ID=%s", id)

	err := s.repo.Delete(ctx, id)

	if err != nil {
		log.Printf("Failed to delete task: ID=%s, Error=%v", id, err)

		return err
	}

	log.Printf("Task deleted successfully: ID=%s", id)

	return nil
}

func (s *TaskService) ListTasks(ctx context.Context, filters map[string]interface{}) ([]models.Task, error) {
	log.Printf("Listing tasks with filters: %+v", filters)
	tasks, err := s.repo.List(ctx, filters)

	if err != nil {
		log.Printf("Failed to list tasks: Error=%v", err)

		return nil, err
	}

	log.Printf("Listed tasks successfully: Found %d tasks", len(tasks))

	return tasks, nil
}

func (s *TaskService) DuplicateTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	log.Printf("Duplicating task: ID=%s", id)

	duplicatedTask, err := s.repo.Duplicate(ctx, id)

	if err != nil {
		log.Printf("Failed to duplicate task: ID=%s, Error=%v", id, err)

		return nil, err
	}

	log.Printf("Task duplicated successfully: OriginalID=%s, DuplicatedID=%s", id, duplicatedTask.ID)

	return duplicatedTask, nil
}
