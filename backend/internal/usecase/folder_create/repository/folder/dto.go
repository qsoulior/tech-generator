package folder_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
	"github.com/samber/lo"
)

type folder struct {
	AuthorID     int64  `db:"author_id"`
	RootAuthorID int64  `db:"root_author_id"`
	UserID       int64  `db:"user_id"`
	Role         string `db:"role"`
}

type folders []folder

func (f folders) toDomain() *domain.Folder {
	if len(f) == 0 {
		return nil
	}

	return &domain.Folder{
		AuthorID:     f[0].AuthorID,
		RootAuthorID: f[0].RootAuthorID,
		Users: lo.Map(f, func(folder folder, _ int) domain.FolderUser {
			return domain.FolderUser{
				ID:   folder.UserID,
				Role: user_domain.Role(folder.Role),
			}
		}),
	}
}
