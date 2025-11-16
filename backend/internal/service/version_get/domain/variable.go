package domain

import variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"

type Variable struct {
	ID          int64
	Name        string
	Type        variable_domain.Type
	Expression  *string
	IsInput     bool
	Constraints []Constraint
}
