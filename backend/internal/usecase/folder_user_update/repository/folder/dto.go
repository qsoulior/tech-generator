package folder_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"

type folder struct {
	AuthorID     int64 `db:"author_id"`
	RootAuthorID int64 `db:"root_author_id"`
}

func (t folder) toDomain() *domain.Folder {
	return &domain.Folder{
		AuthorID:     t.AuthorID,
		RootAuthorID: t.RootAuthorID,
	}
}
