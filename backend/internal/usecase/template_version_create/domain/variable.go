package domain

import (
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

type Variable struct {
	Name        string
	Type        variable_domain.Type
	Expression  string
	Constraints []VariableConstraint
}

type VariableToCreate struct {
	VersionID  int64
	Name       string
	Type       variable_domain.Type
	Expression string
}

type VariableConstraint struct {
	Name       string
	Expression string
	IsActive   bool
}

type VariableConstraintToCreate struct {
	VariableID int64
	Name       string
	Expression string
	IsActive   bool
}
