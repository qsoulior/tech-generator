package user_token_delete_handler

import (
	"context"
	"net/http"

	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) UserTokenDelete(_ context.Context) (*api.UserTokenDeleteNoContent, error) {
	cookie := http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	resp := api.UserTokenDeleteNoContent{
		SetCookie: cookie.String(),
	}

	return &resp, nil
}
