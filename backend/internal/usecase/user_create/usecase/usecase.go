package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

type Usecase struct {
	userRepo       userRepository
	passwordHasher passwordHasher
}

func New(userRepo userRepository, passwordHasher passwordHasher) *Usecase {
	return &Usecase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.UserCreateIn) error {
	if err := in.Validate(); err != nil {
		return err
	}

	exists, err := u.userRepo.ExistsByNameOrEmail(ctx, in.Name, in.Email)
	if err != nil {
		return fmt.Errorf("user repo - exists by name or email: %w", err)
	}

	if exists {
		return domain.ErrUserExists
	}

	passwordHash, err := u.passwordHasher.Hash(in.Password)
	if err != nil {
		return fmt.Errorf("password hasher - hash: %w", err)
	}

	err = u.userRepo.Create(ctx, in.Name, in.Email, passwordHash)
	if err != nil {
		return fmt.Errorf("user repo - create: %w", err)
	}

	return nil
}
