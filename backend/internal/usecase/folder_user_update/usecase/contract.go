//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type folderRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Folder, error)
}

type folderUserRepository interface {
	GetByFolderID(ctx context.Context, folderID int64) ([]domain.FolderUser, error)
	Upsert(ctx context.Context, folderID int64, users []domain.FolderUser) error
	Delete(ctx context.Context, folderID int64, userIDs []int64) error
}

type userRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]int64, error)
}
