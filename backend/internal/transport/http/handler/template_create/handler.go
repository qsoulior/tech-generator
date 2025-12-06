package template_create_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TemplateCreate(ctx context.Context, req *api.TemplateCreateRequest, params api.TemplateCreateParams) (api.TemplateCreateRes, error) {
	in := domain.TemplateCreateIn{
		Name:      req.Name,
		ProjectID: req.ProjectID,
		AuthorID:  params.XUserID,
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

		return nil, fmt.Errorf("template create usecase: %w", err)
	}

	return &api.TemplateCreateCreated{}, nil
}
