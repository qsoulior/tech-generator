package template_delete_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/repository/project"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	return usecase.New(projectRepo)
}
