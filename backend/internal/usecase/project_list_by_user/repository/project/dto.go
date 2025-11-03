package project_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"

type project struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (p *project) toDomain() domain.Project {
	return domain.Project{
		ID:   p.ID,
		Name: p.Name,
	}
}
