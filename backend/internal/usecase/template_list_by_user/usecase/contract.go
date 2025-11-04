//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

type templateRepository interface {
	ListByUserID(ctx context.Context, in domain.TemplateListByUserIn) ([]domain.Template, error)
	GetTotalByUserID(ctx context.Context, in domain.TemplateListByUserIn) (int64, error)
}
