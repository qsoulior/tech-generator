package template_create_from_default_usecase

import (
	"github.com/jmoiron/sqlx"

	version_create_service "github.com/qsoulior/tech-generator/backend/internal/service/version_create"
	version_get_service "github.com/qsoulior/tech-generator/backend/internal/service/version_get"
	new_template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/repository/new_template"
	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/repository/project"
	source_template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/repository/source_template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	sourceTemplateRepo := source_template_repository.New(db)
	newTemplateRepo := new_template_repository.New(db)
	versionGetService := version_get_service.New(db)
	versionCreateService := version_create_service.New(db)
	return usecase.New(projectRepo, sourceTemplateRepo, newTemplateRepo, versionGetService, versionCreateService)
}
