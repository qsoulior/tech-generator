package user_repository

import (
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type user struct {
	ID       int64  `db:"id"`
	Password []byte `db:"password"`
}

func (u *user) toDomain() *domain.User {
	return &domain.User{
		ID:       u.ID,
		Password: u.Password,
	}
}
