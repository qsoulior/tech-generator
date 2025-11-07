package domain

import error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"

var (
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

type ProjectDeleteIn struct {
	ProjectID int64
	UserID    int64
}
