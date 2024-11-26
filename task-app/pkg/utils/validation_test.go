package utils

import (
	"fmt"
	"testing"
	"time"

	"task-app/internal/models"
)

func TestValidateTask(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name      string
		task      *models.Task
		expectErr bool
	}{
		{
			name: "Valid Task",
			task: &models.Task{
				Title:       "Complete Project",
				Description: "Finish the task management system",
				Category:    "Development",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityHigh,
				Status:      models.StatusToDo,
			},
			expectErr: false,
		},
		{
			name: "Empty Title",
			task: &models.Task{
				Title:       "",
				Description: "Test description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      models.StatusInProgress,
			},
			expectErr: true,
		},
		{
			name: "Title Too Long",
			task: &models.Task{
				Title:       "This is an extremely long title that exceeds the maximum allowed length of one hundred characters and should trigger a validation error",
				Description: "Test description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      models.StatusInProgress,
			},
			expectErr: true,
		},
		{
			name: "Past Due Date",
			task: &models.Task{
				Title:       "Past Due Task",
				Description: "Test description",
				DueDate:     time.Now().Add(-24 * time.Hour),
				Priority:    models.PriorityLow,
				Status:      models.StatusToDo,
			},
			expectErr: true,
		},
		{
			name: "Invalid Priority",
			task: &models.Task{
				Title:       "Invalid Priority Task",
				Description: "Test description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    "INVALID_PRIORITY",
				Status:      models.StatusToDo,
			},
			expectErr: true,
		},
		{
			name: "Invalid Status",
			task: &models.Task{
				Title:       "Invalid Status Task",
				Description: "Test description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    models.PriorityMedium,
				Status:      "INVALID_STATUS",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateTask(tc.task)

			if tc.expectErr && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("Did not expect an error, but got: %v", err)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name      string
		title     string
		expectErr bool
	}{
		{"Valid Short Title", "Task Title", false},
		{"Empty Title", "", true},
		{"Very Long Title", "This is an extremely long title that exceeds the maximum allowed length of one hundred characters and should trigger a validation error", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.validateTitle(tc.title)

			if tc.expectErr && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("Did not expect an error, but got: %v", err)
			}
		})
	}
}

func TestValidateDueDate(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name      string
		dueDate   time.Time
		expectErr bool
	}{
		{"Future Date", time.Now().Add(24 * time.Hour), false},
		{"Past Date", time.Now().Add(-24 * time.Hour), true},
		{"Distant Future Date", time.Now().AddDate(6, 0, 0), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.validateDueDate(tc.dueDate)

			if tc.expectErr && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("Did not expect an error, but got: %v", err)
			}
		})
	}
}

func TestBatchValidate(t *testing.T) {
	validator := NewValidator()

	validTask := &models.Task{
		Title:       "Valid Task",
		Description: "A valid task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Priority:    models.PriorityMedium,
		Status:      models.StatusToDo,
	}

	invalidTask := &models.Task{
		Title:       "",
		Description: "An invalid task",
		DueDate:     time.Now().Add(-24 * time.Hour),
		Priority:    models.PriorityHigh,
		Status:      models.StatusInProgress,
	}

	tasks := []*models.Task{validTask, invalidTask}

	validationErrors := validator.BatchValidate(tasks)

	if len(validationErrors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(validationErrors))
	}

	if _, exists := validationErrors[1]; !exists {
		t.Errorf("Expected validation error for the second task")
	}
}

func TestFilterValidationError(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedOutput string
	}{
		{
			name:           "Empty Title Error",
			err:            ErrEmptyTitle,
			expectedOutput: "Please provide a valid task title",
		},
		{
			name:           "Title Too Long Error",
			err:            ErrTitleTooLong,
			expectedOutput: "Task title is too long (max 100 characters)",
		},
		{
			name:           "Invalid Priority Error",
			err:            ErrInvalidPriority,
			expectedOutput: "Invalid priority level",
		},
		{
			name:           "Unknown Error",
			err:            fmt.Errorf("unknown error"),
			expectedOutput: "Invalid task details",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := FilterValidationError(tc.err)
			if output != tc.expectedOutput {
				t.Errorf("Expected '%s', got '%s'", tc.expectedOutput, output)
			}
		})
	}
}
