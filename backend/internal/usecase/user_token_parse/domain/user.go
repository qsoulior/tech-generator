package domain

import error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"

var ErrTokenInvalid = error_domain.NewBaseError("token is invalid")

type User struct {
	ID int64
}
