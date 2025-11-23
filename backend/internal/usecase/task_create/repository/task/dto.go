package task_repository

import (
	"database/sql/driver"
	"encoding/json"
)

type payload map[string]any

func (p *payload) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}

	return json.Marshal(p)
}
