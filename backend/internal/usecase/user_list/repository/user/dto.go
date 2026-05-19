package user_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"

type user struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (u *user) toDomain() domain.User {
	return domain.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
