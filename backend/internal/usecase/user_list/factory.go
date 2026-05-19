package user_list_usecase

import (
	"github.com/jmoiron/sqlx"

	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	userRepo := user_repository.New(db)
	return usecase.New(userRepo)
}
