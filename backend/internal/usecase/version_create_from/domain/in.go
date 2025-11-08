package domain

import error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"

var (
	ErrTemplateNotFound = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid  = error_domain.NewBaseError("template is invalid")
	ErrVersionInvalid   = error_domain.NewBaseError("version is invalid")
)

type VersionCreateFromIn struct {
	AuthorID   int64
	TemplateID int64
	VersionID  int64
}
