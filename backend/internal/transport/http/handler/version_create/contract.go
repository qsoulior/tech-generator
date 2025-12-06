//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package version_create_handler

import (
	"context"

	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

type usecase interface {
	Handle(ctx context.Context, in version_create_domain.VersionCreateIn) error
}
