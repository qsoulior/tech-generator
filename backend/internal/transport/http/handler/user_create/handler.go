package user_create_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) UserCreate(ctx context.Context, req *api.UserCreateRequest) (api.UserCreateRes, error) {
	in := domain.UserCreateIn{
		Name:     req.Name,
		Email:    req.Email,
		Password: domain.Password(req.Password),
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("user create usecase: %w", err)
	}

	return &api.UserCreateCreated{}, nil
}
