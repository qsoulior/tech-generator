package template_import_usecase

import (
	"github.com/jmoiron/sqlx"

	version_create_service "github.com/qsoulior/tech-generator/backend/internal/service/version_create"
	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/repository/project"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	templateRepo := template_repository.New(db)
	versionCreateService := version_create_service.New(db)
	return usecase.New(projectRepo, templateRepo, versionCreateService)
}
