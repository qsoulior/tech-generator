package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

type Usecase struct {
	userRepo userRepository
}

func New(userRepo userRepository) *Usecase {
	return &Usecase{
		userRepo: userRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.UserListIn) (*domain.UserListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	users, err := u.userRepo.List(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("user repo - list: %w", err)
	}

	totalUsers, err := u.userRepo.GetTotal(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("user repo - get total: %w", err)
	}

	out := domain.UserListOut{
		Users:      users,
		TotalUsers: totalUsers,
		TotalPages: (totalUsers + in.Size - 1) / in.Size,
	}

	return &out, nil
}
