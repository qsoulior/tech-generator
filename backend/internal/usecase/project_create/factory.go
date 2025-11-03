package project_create_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/repository/project"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	return usecase.New(projectRepo)
}
