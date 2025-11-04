package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
)

var ErrValueInvalid = errors.New("value is invalid")

type ProjectListByUserIn struct {
	Page    int64
	Size    int64
	Filter  ProjectListByUserFilter
	Sorting *sorting_domain.Sorting
}

type ProjectListByUserFilter struct {
	UserID      int64
	ProjectName *string
}

func (in ProjectListByUserIn) Validate() error {
	if in.Page < 1 {
		return error_domain.NewValidationError("page", ErrValueInvalid)
	}

	if in.Size < 1 {
		return error_domain.NewValidationError("size", ErrValueInvalid)
	}

	if in.Sorting != nil {
		if !in.Sorting.Direction.Valid() {
			return error_domain.NewValidationError("sorting.direction", ErrValueInvalid)
		}
	}

	return nil
}
