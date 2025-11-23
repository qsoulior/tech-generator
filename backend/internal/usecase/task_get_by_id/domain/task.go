package domain

import (
	"time"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
)

var (
	ErrTaskNotFound = error_domain.NewBaseError("task not found")
	ErrTaskInvalid  = error_domain.NewBaseError("task is invalid")
)

type Task struct {
	ID          int64
	VersionID   int64
	Status      task_domain.Status
	Payload     map[string]any
	ResultID    *int64
	Error       *task_domain.ProcessError
	CreatorName string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
