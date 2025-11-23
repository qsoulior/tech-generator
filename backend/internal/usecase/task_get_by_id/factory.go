package task_get_by_id_usecase

import (
	"github.com/jmoiron/sqlx"

	result_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/repository/result"
	task_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/repository/task"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	taskRepo := task_repository.New(db)
	versionRepo := version_repository.New(db)
	resultRepo := result_repository.New(db)
	return usecase.New(taskRepo, versionRepo, resultRepo)
}
