//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_delete_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateDeleteIn) error
}
