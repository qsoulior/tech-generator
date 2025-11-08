//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
}

type templateVersionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.TemplateVersion, error)
}

type variableRepository interface {
	ListByVersionID(ctx context.Context, versionID int64) ([]domain.Variable, error)
}

type variableConstraintRepository interface {
	ListByVariableIDs(ctx context.Context, variableIDs []int64) ([]domain.VariableConstraint, error)
}
