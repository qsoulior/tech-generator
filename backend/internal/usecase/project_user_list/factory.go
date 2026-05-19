package project_user_list_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/repository/project"
	project_user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/repository/project_user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	projectUserRepo := project_user_repository.New(db)
	return usecase.New(projectRepo, projectUserRepo)
}
