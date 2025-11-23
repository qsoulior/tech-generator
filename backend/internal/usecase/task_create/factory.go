package task_create_usecase

import (
	"github.com/jmoiron/sqlx"

	task_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/repository/task"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	versionRepo := version_repository.New(db)
	taskRepo := task_repository.New(db)
	return usecase.New(versionRepo, taskRepo)
}
