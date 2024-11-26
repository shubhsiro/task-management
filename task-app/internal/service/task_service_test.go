package service

import (
	"context"
	"testing"
	"time"

	"task-app/internal/models"
	"task-app/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	service := NewTaskService(repo)
	ctx := context.Background()

	tests := []struct {
		name       string
		req        models.CreateTaskRequest
		expectErr  bool
		expectTask *models.Task
	}{
		{
			name: "Valid Task",
			req: models.CreateTaskRequest{
				Title:       "Valid Task",
				Description: "Task Description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      models.StatusToDo,
			},
			expectErr:  false,
			expectTask: &models.Task{Title: "Valid Task"},
		},
		{
			name: "Invalid Task - Empty Title",
			req: models.CreateTaskRequest{
				Title:       "",
				Description: "No Title",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      models.StatusToDo,
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task, err := service.CreateTask(ctx, test.req)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectTask.Title, task.Title)
				assert.NotEqual(t, uuid.Nil, task.ID)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	service := NewTaskService(repo)
	ctx := context.Background()

	// Seed a task
	task := &models.Task{
		Title:       "Seeded Task",
		Description: "Seeded Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name       string
		taskID     uuid.UUID
		expectErr  bool
		expectTask *models.Task
	}{
		{
			name:       "Get Existing Task",
			taskID:     task.ID,
			expectErr:  false,
			expectTask: task,
		},
		{
			name:      "Get Non-existent Task",
			taskID:    uuid.New(),
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			retrievedTask, err := service.GetTask(ctx, test.taskID)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectTask.Title, retrievedTask.Title)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	service := NewTaskService(repo)
	ctx := context.Background()

	// Create a task to update
	task := &models.Task{
		Title:       "Task to Update",
		Description: "This task will be updated",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium, // Ensure this is valid
		Status:      models.StatusToDo,
	}

	err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Define test cases as a list of structs
	testCases := []struct {
		name          string
		updatedTask   *models.Task
		expectError   bool
		expectedTitle string
	}{
		{
			name:          "Update Existing Task",
			updatedTask:   &models.Task{ID: task.ID, Title: "Updated Title", Description: "This task will be updated", DueDate: time.Now().Add(24 * time.Hour), Priority: models.PriorityMedium, Status: models.StatusToDo},
			expectError:   false,
			expectedTitle: "Updated Title",
		},
		{
			name:          "Update Non-existent Task",
			updatedTask:   &models.Task{ID: uuid.New(), Title: "Non-existent Task", Description: "This task does not exist", DueDate: time.Now().Add(24 * time.Hour), Priority: models.PriorityLow, Status: models.StatusToDo},
			expectError:   true,
			expectedTitle: "",
		},
	}

	// Run each test case using a loop
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Update the task
			err = service.UpdateTask(ctx, tc.updatedTask)

			if tc.expectError {
				assert.Error(t, err, "Expected error for non-existent task")
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// Verify the update for existing task
				updatedTask, _ := repo.GetByID(ctx, tc.updatedTask.ID)
				assert.Equal(t, tc.expectedTitle, updatedTask.Title, "Task title should be updated")
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	service := NewTaskService(repo)
	ctx := context.Background()

	// Seed a task
	task := &models.Task{
		Title:       "Task to Delete",
		Description: "To be deleted",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityLow,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name      string
		taskID    uuid.UUID
		expectErr bool
	}{
		{
			name:      "Delete Existing Task",
			taskID:    task.ID,
			expectErr: false,
		},
		{
			name:      "Delete Non-existent Task",
			taskID:    uuid.New(),
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.DeleteTask(ctx, test.taskID)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := service.GetTask(ctx, test.taskID)
				assert.Error(t, err)
			}
		})
	}
}

func TestDuplicateTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	service := NewTaskService(repo)
	ctx := context.Background()

	// Seed a task
	task := &models.Task{
		Title:       "Task to Duplicate",
		Description: "Duplicated Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium,
		Status:      models.StatusToDo,
	}
	_ = repo.Create(ctx, task)

	tests := []struct {
		name       string
		taskID     uuid.UUID
		expectErr  bool
		expectTask *models.Task
	}{
		{
			name:       "Duplicate Existing Task",
			taskID:     task.ID,
			expectErr:  false,
			expectTask: &models.Task{Title: "Task to Duplicate (Copy)"},
		},
		{
			name:      "Duplicate Non-existent Task",
			taskID:    uuid.New(),
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duplicatedTask, err := service.DuplicateTask(ctx, test.taskID)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectTask.Title, duplicatedTask.Title)
			}
		})
	}
}
