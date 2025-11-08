package template_get_by_id_usecase

import (
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/repository/template"
	template_version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/repository/template_version"
	variable_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/repository/variable"
	variable_constraint_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/repository/variable_constraint"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	templateVersionRepo := template_version_repository.New(db)
	variableRepo := variable_repository.New(db)
	variableConstraintRepo := variable_constraint_repository.New(db)
	return usecase.New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)
}
