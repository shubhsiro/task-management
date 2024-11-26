package utils

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"task-app/internal/models"
)

var (
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrTitleTooLong       = errors.New("title cannot exceed 100 characters")
	ErrInvalidPriority    = errors.New("invalid priority level")
	ErrInvalidStatus      = errors.New("invalid task status")
	ErrDueDateInPast      = errors.New("due date cannot be in the past")
	ErrDescriptionTooLong = errors.New("description cannot exceed 500 characters")
	ErrInvalidCategory    = errors.New("category name is invalid")
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateTask(task *models.Task) error {
	if err := v.validateTitle(task.Title); err != nil {
		return err
	}

	if err := v.validateDescription(task.Description); err != nil {
		return err
	}

	if err := v.validateCategory(task.Category); err != nil {
		return err
	}

	if err := v.validatePriority(task.Priority); err != nil {
		return err
	}

	if err := v.validateStatus(task.Status); err != nil {
		return err
	}

	if err := v.validateDueDate(task.DueDate); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateTitle(title string) error {
	// Check if title is empty
	if title == "" {
		return ErrEmptyTitle
	}

	// Check title length
	if len(title) > 100 {
		return ErrTitleTooLong
	}

	if !containsMeaningfulCharacters(title) {
		return errors.New("title must contain meaningful characters")
	}

	return nil
}

func (v *Validator) validateDescription(description string) error {
	// Optional description length check
	if len(description) > 500 {
		return ErrDescriptionTooLong
	}

	return nil
}

func (v *Validator) validateCategory(category string) error {
	// Optional category validation
	if category != "" {
		if len(category) > 50 {
			return ErrInvalidCategory
		}

		// Ensure category contains only letters, spaces, and hyphens
		for _, char := range category {
			if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '-' {
				return ErrInvalidCategory
			}
		}
	}

	return nil
}

func (v *Validator) validatePriority(priority models.Priority) error {
	validPriorities := map[models.Priority]bool{
		models.PriorityLow:    true,
		models.PriorityMedium: true,
		models.PriorityHigh:   true,
	}

	if !validPriorities[priority] {
		return fmt.Errorf("%w: %s", ErrInvalidPriority, priority)
	}

	return nil
}

func (v *Validator) validateStatus(status models.Status) error {
	validStatuses := map[models.Status]bool{
		models.StatusToDo:       true,
		models.StatusInProgress: true,
		models.StatusDone:       true,
		models.StatusBlocked:    true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	return nil
}

func (v *Validator) validateDueDate(dueDate time.Time) error {
	// Ensure due date is not in the past
	now := time.Now()
	if dueDate.Before(now) {
		return ErrDueDateInPast
	}

	// Add maximum future date limit if needed
	maxFutureDate := now.AddDate(5, 0, 0) // 5 years from now
	if dueDate.After(maxFutureDate) {
		return errors.New("due date cannot be more than 5 years in the future")
	}

	return nil
}

// containsMeaningfulCharacters helper function to check for meaningful characters
func containsMeaningfulCharacters(s string) bool {
	letterCount := 0
	for _, char := range s {
		if unicode.IsLetter(char) {
			letterCount++
		}
	}

	return letterCount > 0
}

// FilterValidationError helps to extract more readable error messages
func FilterValidationError(err error) string {
	switch {
	case errors.Is(err, ErrEmptyTitle):
		return "Please provide a valid task title"
	case errors.Is(err, ErrTitleTooLong):
		return "Task title is too long (max 100 characters)"
	case errors.Is(err, ErrInvalidPriority):
		return "Invalid priority level"
	case errors.Is(err, ErrInvalidStatus):
		return "Invalid task status"
	case errors.Is(err, ErrDueDateInPast):
		return "Due date cannot be in the past"
	default:
		return "Invalid task details"
	}
}

// BatchValidate allows validation of multiple tasks
func (v *Validator) BatchValidate(tasks []*models.Task) map[int]error {
	validationErrors := make(map[int]error)

	for i, task := range tasks {
		if err := v.ValidateTask(task); err != nil {
			validationErrors[i] = err
		}
	}

	return validationErrors
}
