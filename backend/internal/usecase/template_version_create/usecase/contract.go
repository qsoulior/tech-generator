//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
	UpdateByID(ctx context.Context, template domain.TemplateToUpdate) error
}

type templateVersionRepository interface {
	Create(ctx context.Context, templateVersion domain.TemplateVersion) (int64, error)
}

type variableRepository interface {
	Create(ctx context.Context, variables []domain.VariableToCreate) ([]int64, error)
}

type variableConstraintRepository interface {
	Create(ctx context.Context, constraints []domain.VariableConstraintToCreate) error
}
