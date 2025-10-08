package domain

import user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"

type Folder struct {
	AuthorID     int64
	RootAuthorID int64
	Users        []FolderUser
}

type FolderUser struct {
	ID   int64
	Role user_domain.Role
}

type FolderToCreate struct {
	ParentID     *int64
	Name         string
	AuthorID     int64
	RootAuthorID int64
}
