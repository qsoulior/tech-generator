package project_list_by_user_usecase

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	return usecase.New(templateRepo)
}
