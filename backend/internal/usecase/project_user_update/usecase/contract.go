//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type projectRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Project, error)
}

type projectUserRepository interface {
	GetByProjectID(ctx context.Context, projectID int64) ([]domain.ProjectUser, error)
	Upsert(ctx context.Context, projectID int64, users []domain.ProjectUser) error
	Delete(ctx context.Context, projectID int64, userIDs []int64) error
}

type userRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]int64, error)
}
