package template_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"

type template struct {
	AuthorID     int64 `db:"author_id"`
	RootAuthorID int64 `db:"root_author_id"`
}

func (t template) toDomain() *domain.Template {
	return &domain.Template{
		AuthorID:     t.AuthorID,
		RootAuthorID: t.RootAuthorID,
	}
}
