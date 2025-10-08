package folder_create_usecase

import (
	"github.com/jmoiron/sqlx"

	folder_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/repository/folder"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	folderRepo := folder_repository.New(db)
	return usecase.New(folderRepo)
}
