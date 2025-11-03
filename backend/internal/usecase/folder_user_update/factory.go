package folder_user_update_usecase

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"

	folder_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/repository/folder"
	folder_user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/repository/folder_user"
	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	folderRepo := folder_repository.New(db)
	folderUserRepo := folder_user_repository.New(db, trmsqlx.DefaultCtxGetter)
	userRepo := user_repository.New(db)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return usecase.New(folderRepo, folderUserRepo, userRepo, trManager)
}
