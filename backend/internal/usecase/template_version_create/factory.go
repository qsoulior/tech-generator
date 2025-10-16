package template_version_create_usecase

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/repository/template"
	template_version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/repository/template_version"
	variable_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/repository/variable"
	variable_constraint_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/repository/variable_constraint"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db, trmsqlx.DefaultCtxGetter)
	templateVersionRepo := template_version_repository.New(db, trmsqlx.DefaultCtxGetter)
	variableRepo := variable_repository.New(db, trmsqlx.DefaultCtxGetter)
	variableConstraintRepo := variable_constraint_repository.New(db, trmsqlx.DefaultCtxGetter)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return usecase.New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo, trManager)
}
