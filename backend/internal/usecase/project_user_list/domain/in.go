package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

type ProjectUserListIn struct {
	UserID    int64
	ProjectID int64
}
