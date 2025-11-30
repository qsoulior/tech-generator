//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package project_create_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.ProjectCreateIn) error
}
