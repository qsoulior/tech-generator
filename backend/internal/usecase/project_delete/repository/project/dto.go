package project_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"

type project struct {
	AuthorID int64 `db:"author_id"`
}

func (t project) toDomain() *domain.Project {
	return &domain.Project{
		AuthorID: t.AuthorID,
	}
}
