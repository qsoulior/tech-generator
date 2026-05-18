package template_get_meta_by_id_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_meta_by_id/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateGetMetaByID(ctx context.Context, params api.TemplateGetMetaByIDParams) (api.TemplateGetMetaByIDRes, error) {
	in := domain.TemplateGetMetaByIDIn{
		TemplateID: params.TemplateID,
		UserID:     params.XUserID,
	}

	out, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("template get meta by id usecase: %w", err)
	}

	return &api.TemplateGetMetaByIDResponse{Name: out.Name}, nil
}
