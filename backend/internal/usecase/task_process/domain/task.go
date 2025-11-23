package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
)

var ErrTaskNotFound = error_domain.NewBaseError("task not found")

type Task struct {
	VersionID int64
	Payload   map[string]any
}

type TaskUpdate struct {
	ID       int64
	Status   task_domain.Status
	ResultID *int64
	Error    *task_domain.ProcessError
}
