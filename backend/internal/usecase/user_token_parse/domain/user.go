package domain

import base_domain "github.com/qsoulior/tech-generator/backend/internal/domain"

var ErrTokenInvalid = base_domain.NewError("token is invalid")

type User struct {
	ID int64
}
