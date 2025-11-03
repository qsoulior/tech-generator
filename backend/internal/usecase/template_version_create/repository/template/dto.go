package template_repository

import (
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type template struct {
	AuthorID        int64  `db:"author_id"`
	ProjectAuthorID int64  `db:"project_author_id"`
	UserID          int64  `db:"user_id"`
	Role            string `db:"role"`
}

type templates []template

func (ts templates) toDomain() *domain.Template {
	if len(ts) == 0 {
		return nil
	}

	return &domain.Template{
		AuthorID:        ts[0].AuthorID,
		ProjectAuthorID: ts[0].ProjectAuthorID,
		Users: lo.Map(ts, func(t template, _ int) domain.TemplateUser {
			return domain.TemplateUser{
				ID:   t.UserID,
				Role: user_domain.Role(t.Role),
			}
		}),
	}
}
