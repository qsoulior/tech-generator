package template_create_from_default_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateCreateFromDefault(ctx context.Context, req *api.TemplateCreateFromDefaultRequest, params api.TemplateCreateFromDefaultParams) (api.TemplateCreateFromDefaultRes, error) {
	in := convertRequestToIn(req, params)

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

		return nil, fmt.Errorf("template create from default usecase: %w", err)
	}

	return &api.TemplateCreateFromDefaultResponse{ID: out.ID}, nil
}

func convertRequestToIn(req *api.TemplateCreateFromDefaultRequest, params api.TemplateCreateFromDefaultParams) domain.TemplateCreateFromDefaultIn {
	return domain.TemplateCreateFromDefaultIn{
		AuthorID:         params.XUserID,
		ProjectID:        req.ProjectID,
		SourceTemplateID: req.SourceTemplateID,
		Name:             req.Name,
	}
}
