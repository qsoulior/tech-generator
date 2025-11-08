package variable_constraint_repository

import (
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type variableConstraint struct {
	ID         int64  `db:"id"`
	VariableID int64  `db:"variable_id"`
	Name       string `db:"name"`
	Expression string `db:"expression"`
	IsActive   bool   `db:"is_active"`
}

func (c *variableConstraint) toDomain() domain.VariableConstraint {
	return domain.VariableConstraint{
		ID:         c.ID,
		VariableID: c.VariableID,
		Name:       c.Name,
		Expression: c.Expression,
		IsActive:   c.IsActive,
	}
}
