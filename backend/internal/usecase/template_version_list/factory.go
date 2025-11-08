package template_version_list_repository

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/repository/template"
	template_version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/repository/template_version"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	templateVersionRepo := template_version_repository.New(db)
	return usecase.New(templateRepo, templateVersionRepo)
}
