package project_update_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectUpdateByID(ctx context.Context, req *api.ProjectUpdateRequest, params api.ProjectUpdateByIDParams) (api.ProjectUpdateByIDRes, error) {
	in := domain.ProjectUpdateIn{
		ProjectID: params.ProjectID,
		UserID:    params.XUserID,
		Name:      req.Name,
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
		return nil, fmt.Errorf("project update by id usecase: %w", err)
	}

	return &api.ProjectUpdateByIDNoContent{}, nil
}
