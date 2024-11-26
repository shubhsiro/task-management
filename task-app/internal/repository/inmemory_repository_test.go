package repository

import (
	"context"
	"testing"
	"time"

	"task-app/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostOperations(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	ctx := context.Background()

	tests := []struct {
		name     string
		task     *models.Task
		expected string
	}{
		{
			name: "Create Task",
			task: &models.Task{
				Title:       "Task 1",
				Description: "Description 1",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      models.StatusToDo,
			},
			expected: "Task 1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Create(ctx, test.task)
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, test.task.ID)
			assert.Equal(t, test.expected, test.task.Title)
		})
	}
}

func TestGetOperations(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	ctx := context.Background()

	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name     string
		id       uuid.UUID
		expected *models.Task
		hasError bool
	}{
		{
			name:     "Get Existing Task",
			id:       task.ID,
			expected: task,
			hasError: false,
		},
		{
			name:     "Get Non-existent Task",
			id:       uuid.New(),
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, test.id)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.Title, result.Title)
			}
		})
	}
}

func TestUpdateOperations(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	ctx := context.Background()

	task := &models.Task{
		Title:       "Initial Task",
		Description: "Initial Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name     string
		update   *models.Task
		expected string
		hasError bool
	}{
		{
			name: "Update Existing Task",
			update: &models.Task{
				ID:    task.ID,
				Title: "Updated Task",
			},
			expected: "Updated Task",
			hasError: false,
		},
		{
			name: "Update Non-existent Task",
			update: &models.Task{
				ID:    uuid.New(),
				Title: "Non-existent Task",
			},
			expected: "",
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Update(ctx, test.update)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				updatedTask, _ := repo.GetByID(ctx, test.update.ID)
				assert.Equal(t, test.expected, updatedTask.Title)
			}
		})
	}
}

func TestDeleteOperations(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	ctx := context.Background()

	task := &models.Task{
		Title:       "Task to Delete",
		Description: "Delete Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityLow,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name     string
		id       uuid.UUID
		hasError bool
	}{
		{
			name:     "Delete Existing Task",
			id:       task.ID,
			hasError: false,
		},
		{
			name:     "Delete Non-existent Task",
			id:       uuid.New(),
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Delete(ctx, test.id)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := repo.GetByID(ctx, test.id)
				assert.Error(t, err)
			}
		})
	}
}
