//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package project_get_by_id_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.ProjectGetByIDIn) (*domain.ProjectGetByIDOut, error)
}
