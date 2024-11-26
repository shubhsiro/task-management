package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task-app/internal/models"
	"task-app/internal/service"
	"task-app/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a task")

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	task, err := h.service.CreateTask(r.Context(), req)
	if err != nil {
		errorMsg := utils.FilterValidationError(err)
		log.Printf("Error creating task: %v\n", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)

		return
	}

	log.Printf("Task created successfully: %v\n", task.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get a task")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v\n", vars["id"])
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		log.Printf("Task not found with ID: %v\n", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Task retrieved successfully: %v\n", task.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update a task")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v\n", vars["id"])
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = id
	if err := h.service.UpdateTask(r.Context(), &task); err != nil {
		log.Printf("Error updating task with ID %v: %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task updated successfully: %v\n", id)
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete a task")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v\n", vars["id"])
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		log.Printf("Error deleting task with ID %v: %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task deleted successfully: %v\n", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list tasks")

	filters := make(map[string]interface{})

	// Parse query parameters for filtering
	query := r.URL.Query()
	if title := query.Get("title"); title != "" {
		filters["title"] = title
	}

	if category := query.Get("category"); category != "" {
		filters["category"] = category
	}

	if status := query.Get("status"); status != "" {
		filters["status"] = models.Status(status)
	}

	if priority := query.Get("priority"); priority != "" {
		filters["priority"] = models.Priority(priority)
	}

	if dueDate := query.Get("due_date"); dueDate != "" {
		if parsedDate, err := time.Parse(time.RFC3339, dueDate); err == nil {
			filters["due_date"] = parsedDate
		} else {
			log.Printf("Error parsing due_date: %v\n", err)
		}
	}

	tasks, err := h.service.ListTasks(r.Context(), filters)
	if err != nil {
		log.Printf("Error retrieving tasks: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Tasks retrieved successfully: %d tasks found\n", len(tasks))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) DuplicateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to duplicate a task")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v\n", vars["id"])
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.DuplicateTask(r.Context(), id)
	if err != nil {
		log.Printf("Error duplicating task with ID %v: %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task duplicated successfully: %v\n", task.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
