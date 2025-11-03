package project_user_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type projectUser struct {
	UserID int64  `db:"user_id"`
	Role   string `db:"role"`
}

func (u projectUser) toDomain() domain.ProjectUser {
	return domain.ProjectUser{
		ID:   u.UserID,
		Role: user_domain.Role(u.Role),
	}
}
