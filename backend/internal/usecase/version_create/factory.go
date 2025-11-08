package version_create_usecase

import (
	"github.com/jmoiron/sqlx"

	version_create_service "github.com/qsoulior/tech-generator/backend/internal/service/version_create"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/version_create/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	versionCreateService := version_create_service.New(db)
	return usecase.New(templateRepo, versionCreateService)
}
