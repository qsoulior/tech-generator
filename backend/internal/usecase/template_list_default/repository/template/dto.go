package template_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

var sortingAttributes = map[string]struct{}{
	"name": {},
}

type template struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (t *template) toDomain() domain.Template {
	return domain.Template{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
