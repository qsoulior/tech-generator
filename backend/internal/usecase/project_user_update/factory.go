package project_user_update_usecase

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"

	project_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/repository/project"
	project_user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/repository/project_user"
	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	projectRepo := project_repository.New(db)
	projectUserRepo := project_user_repository.New(db, trmsqlx.DefaultCtxGetter)
	userRepo := user_repository.New(db)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return usecase.New(projectRepo, projectUserRepo, userRepo, trManager)
}
