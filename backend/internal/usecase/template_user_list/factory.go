package template_user_list_usecase

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/repository/template"
	template_user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/repository/template_user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	templateUserRepo := template_user_repository.New(db)
	return usecase.New(templateRepo, templateUserRepo)
}
