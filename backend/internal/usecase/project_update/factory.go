package project_update_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/repository/project"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	return usecase.New(projectRepo)
}
