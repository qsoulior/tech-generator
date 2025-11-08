package constraint_repository

import (
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type сonstraint struct {
	ID         int64  `db:"id"`
	VariableID int64  `db:"variable_id"`
	Name       string `db:"name"`
	Expression string `db:"expression"`
	IsActive   bool   `db:"is_active"`
}

func (c *сonstraint) toDomain() domain.Constraint {
	return domain.Constraint{
		ID:         c.ID,
		VariableID: c.VariableID,
		Name:       c.Name,
		Expression: c.Expression,
		IsActive:   c.IsActive,
	}
}
