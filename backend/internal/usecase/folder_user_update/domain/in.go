package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrFolderNotFound = error_domain.NewBaseError("folder not found")
	ErrFolderInvalid  = error_domain.NewBaseError("folder is invalid")
)

type FolderUserUpdateIn struct {
	UserID   int64
	FolderID int64
	Users    []FolderUser
}
