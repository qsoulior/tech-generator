//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package version_list_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.VersionListIn) (*domain.VersionListOut, error)
}
