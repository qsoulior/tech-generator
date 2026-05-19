package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var ErrValueInvalid = errors.New("value is invalid")

type UserListIn struct {
	Page   int64
	Size   int64
	Filter UserListFilter
}

type UserListFilter struct {
	ExcludeUserID int64
	UserName      *string
}

func (in UserListIn) Validate() error {
	if in.Page < 1 {
		return error_domain.NewValidationError("page", ErrValueInvalid)
	}

	if in.Size < 1 {
		return error_domain.NewValidationError("size", ErrValueInvalid)
	}

	return nil
}
