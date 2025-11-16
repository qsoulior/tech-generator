package task_repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type task struct {
	VersionID int64   `db:"version_id"`
	Payload   payload `db:"payload"`
}

type payload map[string]any

func (p *payload) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

func (t *task) toDomain() *domain.Task {
	return &domain.Task{
		VersionID: t.VersionID,
		Payload:   t.Payload,
	}
}

type taskError domain.ProcessError

func (e *taskError) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}

	return json.Marshal(e)
}
