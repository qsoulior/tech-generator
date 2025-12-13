//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
}

type versionGetService interface {
	Handle(ctx context.Context, versionID int64) (*version_get_domain.Version, error)
}

type versionCreateService interface {
	Handle(ctx context.Context, in version_create_domain.VersionCreateIn) (int64, error)
}
