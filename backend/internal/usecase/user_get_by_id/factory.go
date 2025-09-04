package user_get_by_id_usecase

import (
	"github.com/jmoiron/sqlx"

	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	userRepo := user_repository.New(db)
	return usecase.New(userRepo)
}
