package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

type ProjectUserUpdateIn struct {
	UserID    int64
	ProjectID int64
	Users     []ProjectUser
}
