package project_get_by_id_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectGetByID(ctx context.Context, params api.ProjectGetByIDParams) (api.ProjectGetByIDRes, error) {
	in := domain.ProjectGetByIDIn{
		ProjectID: params.ProjectID,
		UserID:    params.XUserID,
	}

	out, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("project get by id usecase: %w", err)
	}

	return &api.ProjectGetByIDResponse{
		Name:       out.Name,
		AuthorName: out.AuthorName,
	}, nil
}
