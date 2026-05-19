package task_create_usecase

import (
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"

	task_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/repository/task"
	version_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/repository/version"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/service/publisher"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/usecase"
)

func New(db *sqlx.DB, amqp *amqp091.Channel) *usecase.Usecase {
	versionRepo := version_repository.New(db)
	taskRepo := task_repository.New(db)
	publisher := publisher.New(amqp)
	return usecase.New(versionRepo, taskRepo, publisher)
}
