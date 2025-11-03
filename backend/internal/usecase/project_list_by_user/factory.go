package project_list_by_user_usecase

import (
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/repository/project"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	return usecase.New(projectRepo)
}
