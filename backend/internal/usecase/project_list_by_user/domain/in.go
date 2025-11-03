package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var ErrValueInvalid = errors.New("value is invalid")

type ProjectListByUserIn struct {
	UserID int64
	Page   uint64
	Size   uint64
}

func (in ProjectListByUserIn) Validate() error {
	if in.Page < 1 {
		return error_domain.NewValidationError("page", ErrValueInvalid)
	}

	if in.Size < 1 {
		return error_domain.NewValidationError("size", ErrValueInvalid)
	}

	return nil
}
