package template_create_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/repository/project"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	templateRepo := template_repository.New(db)
	return usecase.New(projectRepo, templateRepo)
}
