package user_create_usecase

import (
	"github.com/jmoiron/sqlx"

	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/service/password_hasher"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/usecase"
)

func New(db *sqlx.DB) *usecase.Usecase {
	userRepo := user_repository.New(db)
	passwordHasher := password_hasher.New()
	return usecase.New(userRepo, passwordHasher)
}
