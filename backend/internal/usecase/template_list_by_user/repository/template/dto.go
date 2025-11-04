package template_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

var sortingAttributes = map[string]struct{}{
	"template_name": {},
}

type template struct {
	ID                 int64      `db:"id"`
	TemplateName       string     `db:"template_name"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          *time.Time `db:"updated_at"`
	TemplateAuthorName string     `db:"template_author_name"`
}

func (t *template) toDomain() domain.Template {
	return domain.Template{
		ID:         t.ID,
		Name:       t.TemplateName,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
		AuthorName: t.TemplateAuthorName,
	}
}
