package version_get_service

import (
	"github.com/jmoiron/sqlx"

	constraint_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_get/repository/constraint"
	variable_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_get/repository/variable"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/service/version_get/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/service"
)

func New(db *sqlx.DB) *service.Service {
	versionRepo := version_repository.New(db)
	variableRepo := variable_repository.New(db)
	constraintRepo := constraint_repository.New(db)
	return service.New(versionRepo, variableRepo, constraintRepo)
}
