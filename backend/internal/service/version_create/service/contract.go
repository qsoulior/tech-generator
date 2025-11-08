//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package service

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

type templateRepository interface {
	UpdateByID(ctx context.Context, template domain.TemplateToUpdate) error
}

type versionRepository interface {
	Create(ctx context.Context, version domain.Version) (int64, error)
}

type variableRepository interface {
	Create(ctx context.Context, variables []domain.VariableToCreate) ([]int64, error)
}

type constraintRepository interface {
	Create(ctx context.Context, constraints []domain.ConstraintToCreate) error
}
