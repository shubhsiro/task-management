# Task Management System

## Overview

This is a robust, in-memory task management system built with Go, featuring comprehensive CRUD operations, advanced filtering, and detailed validation.

### Features

- In-memory task storage
- Full CRUD operations
- Task validation
- Filtering and searching tasks
- Flexible task management
- RESTful API endpoints


## Prerequisites

- Go 1.21+
- Git

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/shubhsiro/task-management.git
   
   cd task-management

2. Install dependencies
    ```bash
    go mod tidy

3. Run the application:
    ```bash
    go run cmd/main.go

### API Endpoints
## Tasks
- POST /tasks: Create a new task
- GET /tasks: List tasks (with optional filtering)
- GET /tasks/{id}: Get a specific task
- PUT /tasks/{id}: Update a task
- DELETE /tasks/{id}: Delete a task
- POST /tasks/{id}/duplicate: Duplicate a task

### Filtering Parameters
## You can filter tasks by the following parameters:

- title: Filter by task title
- category: Filter by task category
- status: Filter by task status (TODO, IN_PROGRESS, DONE, BLOCKED)
- priority: Filter by priority (LOW, MEDIUM, HIGH)
- due_date: Filter by due date

### Task Model
A task consists of the following fields:

- ID: Unique ID (UUID)
- Title: Title of the task (max 100 characters)
- Description: Description of the task (max 500 characters)
- Category: Category of the task
- Due Date: Due date of the task (must be in the future, within 5 years)
- Priority: Task priority (LOW, MEDIUM, HIGH)
- Status: Task status (TODO, IN_PROGRESS, DONE, BLOCKED)
- Creation Timestamp: Timestamp when the task was created
- Update Timestamp: Timestamp when the task was last updated

## Validation Rules
- Title: Required, max 100 characters
- Description: Optional, max 500 characters
- Category: Optional, letters, spaces, and hyphens allowed
- Priority: LOW, MEDIUM, HIGH
- Status: TODO, IN_PROGRESS, DONE, BLOCKED
- Due Date: Must be in the future, within 5 years

### Running Tests
To run the unit tests:
```bash
    go test ./...
```

### Contact
Your Name - shubham.sirothiya@gmail.com
Project Link: https://github.com/shubhsiro/task-management