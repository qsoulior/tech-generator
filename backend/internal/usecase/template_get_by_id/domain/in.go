package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrTemplateNotFound        = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid         = error_domain.NewBaseError("template is invalid")
	ErrTemplateVersionNotFound = errors.New("template version not found")
)

type TemplateGetByIDIn struct {
	TemplateID int64
	UserID     int64
}
