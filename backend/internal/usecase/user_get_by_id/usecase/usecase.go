package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type Usecase struct {
	userRepo userRepository
}

func New(userRepo userRepository) *Usecase {
	return &Usecase{
		userRepo: userRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user repo - get by id: %w", err)
	}

	return user, nil
}
