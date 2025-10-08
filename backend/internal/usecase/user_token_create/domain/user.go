package domain

import error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"

var ErrUserDoesNotExist = error_domain.NewBaseError("user does not exist")

type User struct {
	ID       int64
	Password []byte
}
