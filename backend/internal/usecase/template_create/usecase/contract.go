//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

type folderRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Folder, error)
}

type templateRepository interface {
	Create(ctx context.Context, template domain.Template) error
}
