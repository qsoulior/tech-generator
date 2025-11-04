package project_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"

var sortingAttributes = map[string]struct{}{
	"project_name": {},
}

type project struct {
	ID          int64  `db:"id"`
	ProjectName string `db:"project_name"`
	AuthorName  string `db:"author_name"`
}

func (p *project) toDomain() domain.Project {
	return domain.Project{
		ID:         p.ID,
		Name:       p.ProjectName,
		AuthorName: p.AuthorName,
	}
}
