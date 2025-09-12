package user_token_parse_usecase

import (
	"crypto/ed25519"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/service/token_parser"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/usecase"
)

func New(publicKey ed25519.PublicKey) *usecase.Usecase {
	tokenParser := token_parser.New(publicKey)
	return usecase.New(tokenParser)
}
