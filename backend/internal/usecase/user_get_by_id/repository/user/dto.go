package user_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type user struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

func (u *user) toDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
