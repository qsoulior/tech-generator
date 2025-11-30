package task_repository

import (
	"encoding/json"
	"errors"
	"time"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

type task struct {
	ID          int64      `db:"id"`
	VersionID   int64      `db:"version_id"`
	Status      string     `db:"status"`
	Payload     payload    `db:"payload"`
	ResultID    *int64     `db:"result_id"`
	Error       *taskError `db:"error"`
	CreatorName string     `db:"creator_name"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func (t *task) toDomain() *domain.Task {
	return &domain.Task{
		ID:          t.ID,
		VersionID:   t.VersionID,
		Status:      task_domain.Status(t.Status),
		Payload:     t.Payload,
		ResultID:    t.ResultID,
		Error:       (*task_domain.ProcessError)(t.Error),
		CreatorName: t.CreatorName,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

type payload map[string]string

func (p *payload) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

type taskError task_domain.ProcessError

func (e *taskError) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}
