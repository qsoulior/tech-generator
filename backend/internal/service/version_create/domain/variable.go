package domain

import (
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

type Variable struct {
	Name        string
	Type        variable_domain.Type
	Expression  string
	Constraints []Constraint
}

type VariableToCreate struct {
	VersionID  int64
	Name       string
	Type       variable_domain.Type
	Expression string
}
