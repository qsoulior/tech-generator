package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

type Usecase struct {
	tokenParser tokenParser
}

func New(tokenParser tokenParser) *Usecase {
	return &Usecase{
		tokenParser: tokenParser,
	}
}

func (u *Usecase) Handle(_ context.Context, token string) (*domain.User, error) {
	user, err := u.tokenParser.Parse(token)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	return user, nil
}
