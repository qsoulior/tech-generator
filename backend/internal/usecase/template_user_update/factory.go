package template_user_update_usecase

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"

	template_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/repository/template"
	template_user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/repository/template_user"
	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	templateRepo := template_repository.New(db)
	templateUserRepo := template_user_repository.New(db, trmsqlx.DefaultCtxGetter)
	userRepo := user_repository.New(db)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return usecase.New(templateRepo, templateUserRepo, userRepo, trManager)
}
