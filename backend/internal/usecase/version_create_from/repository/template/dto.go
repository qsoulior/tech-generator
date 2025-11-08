package template_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"

type template struct {
	AuthorID        int64 `db:"author_id"`
	ProjectAuthorID int64 `db:"project_author_id"`
}

func (t template) toDomain() *domain.Template {
	return &domain.Template{
		AuthorID:        t.AuthorID,
		ProjectAuthorID: t.ProjectAuthorID,
	}
}
