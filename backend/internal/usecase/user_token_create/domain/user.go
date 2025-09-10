package domain

import base_domain "github.com/qsoulior/tech-generator/backend/internal/domain"

var ErrUserDoesNotExist = base_domain.NewError("user does not exist")

type User struct {
	ID       int64
	Password []byte
}
