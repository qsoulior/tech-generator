package task_process_usecase

import (
	"github.com/jmoiron/sqlx"

	version_get_service "github.com/qsoulior/tech-generator/backend/internal/service/version_get"
	task_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/repository/task"
	data_process_service "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/service/data_process"
	variable_process_service "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/service/variable_process"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	taskRepo := task_repository.New(db)
	versionGetService := version_get_service.New(db)
	variableProcessService := variable_process_service.New()
	dataProcessService := data_process_service.New()
	return usecase.New(taskRepo, versionGetService, variableProcessService, dataProcessService)
}
