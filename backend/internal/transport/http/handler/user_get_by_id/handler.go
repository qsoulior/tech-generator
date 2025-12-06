package user_get_by_id_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) UserGetByID(ctx context.Context, params api.UserGetByIDParams) (api.UserGetByIDRes, error) {
	out, err := h.usecase.Handle(ctx, params.XUserID)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("user get by id usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertOutToResponse(out domain.User) api.UserGetByIDResponse {
	return api.UserGetByIDResponse{
		ID:        out.ID,
		Name:      out.Name,
		Email:     out.Email,
		CreatedAt: out.CreatedAt,
	}
}
