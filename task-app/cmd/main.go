package main

import (
	"log"
	"net/http"

	"task-app/internal/handler"
	"task-app/internal/repository"
	"task-app/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	taskRepo := repository.NewInMemoryTaskRepository()

	taskService := service.NewTaskService(taskRepo)

	taskHandler := handler.NewTaskHandler(taskService)

	router := mux.NewRouter()

	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}/duplicate", taskHandler.DuplicateTask).Methods("POST")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
