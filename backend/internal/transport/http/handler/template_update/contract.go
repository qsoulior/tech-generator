//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_update_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_update/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateUpdateIn) error
}
