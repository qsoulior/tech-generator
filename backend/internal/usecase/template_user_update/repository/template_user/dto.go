package template_user_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

type templateUser struct {
	UserID int64  `db:"user_id"`
	Role   string `db:"role"`
}

func (u templateUser) toDomain() domain.TemplateUser {
	return domain.TemplateUser{
		ID:   u.UserID,
		Role: user_domain.Role(u.Role),
	}
}
