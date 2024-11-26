package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"task-app/internal/models"

	"github.com/google/uuid"
)

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*models.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (r *InMemoryTaskRepository) Create(ctx context.Context, task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task

	log.Printf("Created task: ID=%s, Title=%s, Category=%s, Status=%s", task.ID, task.Title, task.Category, task.Status)

	return nil
}

func (r *InMemoryTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		log.Printf("Task not found: ID=%s", id)

		return nil, fmt.Errorf("task not found")
	}

	log.Printf("Retrieved task: ID=%s, Title=%s, Category=%s", task.ID, task.Title, task.Category)

	return task, nil
}

func (r *InMemoryTaskRepository) Update(ctx context.Context, task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		log.Printf("Task not found for update: ID=%s", task.ID)
		return fmt.Errorf("task not found")
	}

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task

	log.Printf("Updated task: ID=%s, Title=%s, Category=%s, Status=%s", task.ID, task.Title, task.Category, task.Status)

	return nil
}

func (r *InMemoryTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		log.Printf("Task not found for deletion: ID=%s", id)
		return fmt.Errorf("task not found")
	}

	delete(r.tasks, id)

	log.Printf("Deleted task: ID=%s", id)

	return nil
}

func (r *InMemoryTaskRepository) List(ctx context.Context, filters map[string]interface{}) ([]models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filteredTasks []models.Task
	for _, task := range r.tasks {
		if matchesFilters(task, filters) {
			filteredTasks = append(filteredTasks, *task)
		}
	}

	log.Printf("Listed tasks with filters: %+v, Found: %d tasks", filters, len(filteredTasks))

	return filteredTasks, nil
}

func (r *InMemoryTaskRepository) Duplicate(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	originalTask, exists := r.tasks[id]
	if !exists {
		log.Printf("Task not found for duplication: ID=%s", id)
		return nil, fmt.Errorf("task not found")
	}

	duplicatedTask := &models.Task{
		ID:          uuid.New(),
		Title:       originalTask.Title + " (Copy)",
		Description: originalTask.Description,
		Category:    originalTask.Category,
		DueDate:     originalTask.DueDate,
		Priority:    originalTask.Priority,
		Status:      models.StatusToDo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.tasks[duplicatedTask.ID] = duplicatedTask

	log.Printf("Duplicated task: OriginalID=%s, NewID=%s, Title=%s", id, duplicatedTask.ID, duplicatedTask.Title)

	return duplicatedTask, nil
}

func matchesFilters(task *models.Task, filters map[string]interface{}) bool {
	for key, value := range filters {
		switch key {
		case "title":
			if task.Title != value {
				return false
			}
		case "category":
			if task.Category != value {
				return false
			}
		case "status":
			if task.Status != value {
				return false
			}
		case "priority":
			if task.Priority != value {
				return false
			}
		case "due_date":
			dueDate, ok := value.(time.Time)
			if !ok || !task.DueDate.Equal(dueDate) {
				return false
			}
		}
	}
	return true
}
