package domain

import version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"

type TaskProcessIn struct {
	TaskID int64
}

type Version = version_get_domain.Version

type Variable = version_get_domain.Variable

type Constraint = version_get_domain.Constraint

type VariableProcessIn struct {
	Variables []Variable
	Payload   map[string]string
}

type DataProcessIn struct {
	Values map[string]any
	Data   []byte
}
