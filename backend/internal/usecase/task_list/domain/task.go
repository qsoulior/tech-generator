package domain

import (
	"time"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
)

type Task struct {
	ID          int64
	Status      task_domain.Status
	CreatorName string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
