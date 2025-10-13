package folder_create_usecase

import (
	"github.com/jmoiron/sqlx"

	folder_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/repository/folder"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	folderRepo := folder_repository.New(db)
	templateRepo := template_repository.New(db)
	return usecase.New(folderRepo, templateRepo)
}
