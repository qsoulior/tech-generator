package project_get_by_id_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/repository/project"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	return usecase.New(projectRepo)
}
