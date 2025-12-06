package template_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TemplateList(ctx context.Context, params api.TemplateListParams) (api.TemplateListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("template list by user usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.TemplateListParams) domain.TemplateListByUserIn {
	filter := domain.TemplateListByUserFilter{
		UserID:    params.XUserID,
		ProjectID: params.ProjectID,
	}

	if templateName, ok := params.TemplateName.Get(); ok {
		filter.TemplateName = &templateName
	}

	in := domain.TemplateListByUserIn{
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

func convertOutToResponse(out domain.TemplateListByUserOut) api.TemplateListResponse {
	return api.TemplateListResponse{
		Templates: lo.Map(out.Templates, func(t domain.Template, _ int) api.TemplateListResponseTemplatesItem {
			item := api.TemplateListResponseTemplatesItem{
				ID:         t.ID,
				Name:       t.Name,
				AuthorName: t.AuthorName,
				CreatedAt:  t.CreatedAt,
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
