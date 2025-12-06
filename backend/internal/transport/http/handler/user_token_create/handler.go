package user_token_create_handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) UserTokenCreate(ctx context.Context, req *api.UserTokenCreateRequest) (api.UserTokenCreateRes, error) {
	in := domain.UserCreateTokenIn{
		Name:     req.Name,
		Password: domain.Password(req.Password),
	}

	out, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("user token create usecase: %w", err)
	}

	cookie := http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    out.Token,
		Secure:   true,
		HttpOnly: true,
		Expires:  out.ExpiresAt,
	}

	resp := api.UserTokenCreateCreated{
		SetCookie: cookie.String(),
	}

	return &resp, nil
}
