package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type Usecase struct {
	userRepo         userRepository
	passwordVerifier passwordVerifier
	tokenBuilder     tokenBuilder
	tokenExpiration  time.Duration
}

func New(userRepo userRepository, passwordVerifier passwordVerifier, tokenBuilder tokenBuilder, tokenExpiration time.Duration) *Usecase {
	return &Usecase{
		userRepo:         userRepo,
		passwordVerifier: passwordVerifier,
		tokenBuilder:     tokenBuilder,
		tokenExpiration:  tokenExpiration,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.UserCreateTokenIn) (domain.UserCreateTokenOut, error) {
	if err := in.Validate(); err != nil {
		return domain.UserCreateTokenOut{}, err
	}

	user, err := u.userRepo.GetByName(ctx, in.Name)
	if err != nil {
		return domain.UserCreateTokenOut{}, fmt.Errorf("user repo - get by name: %w", err)
	}

	if user == nil {
		return domain.UserCreateTokenOut{}, domain.ErrUserDoesNotExist
	}

	err = u.passwordVerifier.Verify(user.Password, in.Password)
	if err != nil {
		return domain.UserCreateTokenOut{}, domain.ErrPasswordIncorrect
	}

	token, err := u.tokenBuilder.Build(user.ID, u.tokenExpiration)
	if err != nil {
		return domain.UserCreateTokenOut{}, fmt.Errorf("token builder - build: %w", err)
	}

	return token, nil
}
