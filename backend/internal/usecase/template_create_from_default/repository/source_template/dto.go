package source_template_repository

import (
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/domain"
)

type sourceTemplate struct {
	ID            int64  `db:"id"`
	IsDefault     bool   `db:"is_default"`
	LastVersionID *int64 `db:"last_version_id"`
}

func (t *sourceTemplate) toDomain() *domain.SourceTemplate {
	return &domain.SourceTemplate{
		ID:            t.ID,
		IsDefault:     t.IsDefault,
		LastVersionID: t.LastVersionID,
	}
}
