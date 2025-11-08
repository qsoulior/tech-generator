//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
}

type templateVersionRepository interface {
	ListByTemplateID(ctx context.Context, templateID int64) ([]domain.TemplateVersion, error)
}
