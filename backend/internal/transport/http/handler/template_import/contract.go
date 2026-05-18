//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_import_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateImportIn) (*domain.TemplateImportOut, error)
}
