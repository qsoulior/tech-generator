package template_version_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/domain"
)

type templateVersion struct {
	ID         int64     `db:"id"`
	Number     int64     `db:"number"`
	AuthorName string    `db:"author_name"`
	CreatedAt  time.Time `db:"created_at"`
}

func (v *templateVersion) toDomain() domain.TemplateVersion {
	return domain.TemplateVersion{
		ID:         v.ID,
		Number:     v.Number,
		AuthorName: v.AuthorName,
		CreatedAt:  v.CreatedAt,
	}
}
