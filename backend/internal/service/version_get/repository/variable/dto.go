package variable_repository

import (
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type variable struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Type       string `db:"type"`
	Expression string `db:"expression"`
}

func (v *variable) toDomain() domain.Variable {
	return domain.Variable{
		ID:         v.ID,
		Name:       v.Name,
		Type:       variable_domain.Type(v.Type),
		Expression: v.Expression,
	}
}
