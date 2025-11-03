package folder_user_repository

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type folderUser struct {
	UserID int64  `db:"user_id"`
	Role   string `db:"role"`
}

func (u folderUser) toDomain() domain.FolderUser {
	return domain.FolderUser{
		ID:   u.UserID,
		Role: user_domain.Role(u.Role),
	}
}
