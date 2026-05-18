//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_default_list_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateListDefaultIn) (*domain.TemplateListDefaultOut, error)
}
