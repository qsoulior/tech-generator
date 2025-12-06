//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_list_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateListByUserIn) (*domain.TemplateListByUserOut, error)
}
