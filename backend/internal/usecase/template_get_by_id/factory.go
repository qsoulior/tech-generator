package template_get_by_id_usecase

import (
	"github.com/jmoiron/sqlx"

	version_get_service "github.com/qsoulior/tech-generator/backend/internal/service/version_get"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	versionGetService := version_get_service.New(db)
	return usecase.New(templateRepo, versionGetService)
}
