package template_user_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

type templateUser struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}

func (u templateUser) toDomain() domain.TemplateUser {
	return domain.TemplateUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  user_domain.Role(u.Role),
	}
}
