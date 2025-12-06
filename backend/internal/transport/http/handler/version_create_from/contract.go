//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package version_create_from_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.VersionCreateFromIn) error
}
