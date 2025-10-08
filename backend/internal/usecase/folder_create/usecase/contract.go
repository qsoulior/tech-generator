//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
)

type folderRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Folder, error)
	Create(ctx context.Context, name string, authorID int64, rootAuthorID int64, parentID *int64) error
}
