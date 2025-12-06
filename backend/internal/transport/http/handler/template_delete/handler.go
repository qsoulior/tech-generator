package template_delete_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TemplateDeleteByID(ctx context.Context, params api.TemplateDeleteByIDParams) (api.TemplateDeleteByIDRes, error) {
	in := domain.TemplateDeleteIn{
		TemplateID: params.TemplateID,
		UserID:     params.XUserID,
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("template delete by id usecase: %w", err)
	}

	return &api.TemplateDeleteByIDNoContent{}, nil
}
