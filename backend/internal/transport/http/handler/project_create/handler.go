package project_create_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectCreate(ctx context.Context, req *api.ProjectCreateRequest, params api.ProjectCreateParams) (api.ProjectCreateRes, error) {
	in := domain.ProjectCreateIn{
		Name:     req.Name,
		AuthorID: params.XUserID,
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("project create usecase: %w", err)
	}

	return &api.ProjectCreateCreated{}, nil
}
