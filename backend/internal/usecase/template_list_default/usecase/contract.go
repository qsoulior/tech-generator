//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type templateRepository interface {
	ListDefault(ctx context.Context, in domain.TemplateListDefaultIn) ([]domain.Template, error)
	GetTotalDefault(ctx context.Context, in domain.TemplateListDefaultIn) (int64, error)
}
