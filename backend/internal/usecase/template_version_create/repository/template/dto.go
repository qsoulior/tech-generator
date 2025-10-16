package template_repository

import (
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type template struct {
	AuthorID     int64  `db:"author_id"`
	RootAuthorID int64  `db:"root_author_id"`
	UserID       int64  `db:"user_id"`
	Role         string `db:"role"`
}

type templates []template

func (f templates) toDomain() *domain.Template {
	if len(f) == 0 {
		return nil
	}

	return &domain.Template{
		AuthorID:     f[0].AuthorID,
		RootAuthorID: f[0].RootAuthorID,
		Users: lo.Map(f, func(folder template, _ int) domain.TemplateUser {
			return domain.TemplateUser{
				ID:   folder.UserID,
				Role: user_domain.Role(folder.Role),
			}
		}),
	}
}
