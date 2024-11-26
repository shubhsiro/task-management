package models

import (
	"github.com/google/uuid"
	"time"
)

type Priority string
type Status string

const (
	// Priorities
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"

	// Statuses
	StatusToDo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
	StatusBlocked    Status = "BLOCKED"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	DueDate     time.Time `json:"due_date"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	DueDate     time.Time `json:"due_date"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
}
