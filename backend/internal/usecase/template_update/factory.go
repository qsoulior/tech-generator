package template_update_usecase

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_update/repository/template"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_update/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	return usecase.New(templateRepo)
}
