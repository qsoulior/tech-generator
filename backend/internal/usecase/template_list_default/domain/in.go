package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
)

var ErrValueInvalid = errors.New("value is invalid")

type TemplateListDefaultIn struct {
	Page    int64
	Size    int64
	Filter  TemplateListDefaultFilter
	Sorting *sorting_domain.Sorting
}

type TemplateListDefaultFilter struct {
	TemplateName *string
}

func (in TemplateListDefaultIn) Validate() error {
	if in.Page < 1 {
		return error_domain.NewValidationError("page", ErrValueInvalid)
	}

	if in.Size < 1 {
		return error_domain.NewValidationError("size", ErrValueInvalid)
	}

	if in.Sorting != nil && !in.Sorting.Direction.Valid() {
		return error_domain.NewValidationError("sorting.direction", ErrValueInvalid)
	}

	return nil
}
