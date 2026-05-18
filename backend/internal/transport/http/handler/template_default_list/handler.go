package template_default_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateDefaultList(ctx context.Context, params api.TemplateDefaultListParams) (api.TemplateDefaultListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("template list default usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.TemplateDefaultListParams) domain.TemplateListDefaultIn {
	filter := domain.TemplateListDefaultFilter{}

	if templateName, ok := params.TemplateName.Get(); ok {
		filter.TemplateName = &templateName
	}

	in := domain.TemplateListDefaultIn{
		Page:   params.Page,
		Size:   params.Size,
		Filter: filter,
	}

	if sorting, ok := params.Sorting.Get(); ok {
		in.Sorting = &sorting_domain.Sorting{
			Attribute: sorting.Attribute,
			Direction: sorting_domain.SortingDirection(sorting.Direction),
		}
	}

	return in
}

func convertOutToResponse(out domain.TemplateListDefaultOut) api.TemplateDefaultListResponse {
	return api.TemplateDefaultListResponse{
		Templates: lo.Map(out.Templates, func(t domain.Template, _ int) api.TemplateDefaultListResponseTemplatesItem {
			item := api.TemplateDefaultListResponseTemplatesItem{
				ID:        t.ID,
				Name:      t.Name,
				CreatedAt: t.CreatedAt,
			}

			if t.UpdatedAt != nil {
				item.UpdatedAt.SetTo(*t.UpdatedAt)
			}

			return item
		}),
		TotalTemplates: out.TotalTemplates,
		TotalPages:     out.TotalPages,
	}
}
