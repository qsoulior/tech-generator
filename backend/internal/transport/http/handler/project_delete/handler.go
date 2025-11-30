package project_delete_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) ProjectDeleteByID(ctx context.Context, params api.ProjectDeleteByIDParams) (api.ProjectDeleteByIDRes, error) {
	in := domain.ProjectDeleteIn{
		ProjectID: params.ProjectID,
		UserID:    params.XUserID,
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("project delete by id usecase: %w", err)
	}

	return &api.ProjectDeleteByIDNoContent{}, nil
}
