//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
}

type templateUserRepository interface {
	GetByTemplateID(ctx context.Context, templateID int64) ([]domain.TemplateUser, error)
	Upsert(ctx context.Context, templateID int64, users []domain.TemplateUser) error
	Delete(ctx context.Context, templateID int64, userIDs []int64) error
}

type userRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]int64, error)
}
