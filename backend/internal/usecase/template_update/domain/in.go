package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrTemplateNotFound = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid  = error_domain.NewBaseError("template is invalid")
	ErrValueEmpty       = errors.New("value is empty")
)

type TemplateUpdateIn struct {
	TemplateID int64
	UserID     int64
	Name       string
}

func (in TemplateUpdateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
