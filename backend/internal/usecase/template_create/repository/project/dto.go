package project_repository

import (
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

type project struct {
	AuthorID int64  `db:"author_id"`
	UserID   int64  `db:"user_id"`
	Role     string `db:"role"`
}

type projects []project

func (ps projects) toDomain() *domain.Project {
	if len(ps) == 0 {
		return nil
	}

	return &domain.Project{
		AuthorID: ps[0].AuthorID,
		Users: lo.Map(ps, func(p project, _ int) domain.ProjectUser {
			return domain.ProjectUser{
				ID:   p.UserID,
				Role: user_domain.Role(p.Role),
			}
		}),
	}
}
