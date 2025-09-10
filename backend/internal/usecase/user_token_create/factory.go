package user_token_create_usecase

import (
	"crypto/ed25519"

	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/config"
	user_repository "github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/repository/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/service/password_verifier"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/service/token_builder"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/usecase"
)

func New(db *sqlx.DB, privateKey ed25519.PrivateKey, cfg *config.Config) *usecase.Usecase {
	userRepo := user_repository.New(db)
	passwordVerifier := password_verifier.New()
	tokenBuilder := token_builder.New(privateKey)

	return usecase.New(userRepo, passwordVerifier, tokenBuilder, cfg.UserTokenExpiration)
}
