package template_version_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type templateVersion struct {
	ID        int64     `db:"id"`
	Number    int64     `db:"number"`
	CreatedAt time.Time `db:"created_at"`
	Data      []byte    `db:"data"`
}

func (v *templateVersion) toDomain() *domain.TemplateVersion {
	return &domain.TemplateVersion{
		ID:        v.ID,
		Number:    v.Number,
		CreatedAt: v.CreatedAt,
		Data:      v.Data,
	}
}
