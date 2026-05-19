package project_user_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type projectUser struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}

func (u projectUser) toDomain() domain.ProjectUser {
	return domain.ProjectUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  user_domain.Role(u.Role),
	}
}
