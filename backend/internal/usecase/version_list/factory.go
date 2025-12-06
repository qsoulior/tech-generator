package version_list_usecase

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/repository/template"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	versionRepo := version_repository.New(db)
	return usecase.New(templateRepo, versionRepo)
}
