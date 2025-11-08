package version_create_service

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"

	constraint_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_create/repository/constraint"
	template_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_create/repository/template"
	variable_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_create/repository/variable"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_create/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/service"
)

func New(db *sqlx.DB) *service.Service {
	templateRepo := template_repository.New(db, trmsqlx.DefaultCtxGetter)
	versionRepo := version_repository.New(db, trmsqlx.DefaultCtxGetter)
	variableRepo := variable_repository.New(db, trmsqlx.DefaultCtxGetter)
	constraintRepo := constraint_repository.New(db, trmsqlx.DefaultCtxGetter)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return service.New(templateRepo, versionRepo, variableRepo, constraintRepo, trManager)
}
