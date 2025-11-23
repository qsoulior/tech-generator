package task_repository

import (
	"time"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

var sortingAttributes = map[string]struct{}{
	"created_at": {},
	"updated_at": {},
}

type task struct {
	ID          int64      `db:"id"`
	Status      string     `db:"status"`
	CreatorName string     `db:"creator_name"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func (t *task) toDomain() domain.Task {
	return domain.Task{
		ID:          t.ID,
		Status:      task_domain.Status(t.Status),
		CreatorName: t.CreatorName,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
