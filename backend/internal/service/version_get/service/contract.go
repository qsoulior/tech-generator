//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package service

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type versionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Version, error)
}

type variableRepository interface {
	ListByVersionID(ctx context.Context, versionID int64) ([]domain.Variable, error)
}

type constraintRepository interface {
	ListByVariableIDs(ctx context.Context, variableIDs []int64) ([]domain.Constraint, error)
}
