package template_get_by_id_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TemplateGetByID(ctx context.Context, params api.TemplateGetByIDParams) (api.TemplateGetByIDRes, error) {
	in := domain.TemplateGetByIDIn{
		TemplateID: params.TemplateID,
		UserID:     params.XUserID,
	}

	out, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("template get by id usecase: %w", err)
	}

	var resp api.TemplateGetByIDResponse
	if out.Version != nil {
		resp.Version.SetTo(convertVersionToResponse(*out.Version))
	}

	return &resp, nil
}

func convertVersionToResponse(version version_get_domain.Version) api.TemplateGetByIDVersion {
	return api.TemplateGetByIDVersion{
		Number:    version.Number,
		CreatedAt: version.CreatedAt,
		Data:      version.Data,
		Variables: convertVariablesToResponse(version.Variables),
	}
}

func convertVariablesToResponse(variables []version_get_domain.Variable) []api.TemplateGetByIDVersionVariablesItem {
	return lo.Map(variables, func(v version_get_domain.Variable, _ int) api.TemplateGetByIDVersionVariablesItem {
		item := api.TemplateGetByIDVersionVariablesItem{
			ID:          v.ID,
			Name:        v.Name,
			Type:        api.TemplateGetByIDVersionVariablesItemType(v.Type),
			IsInput:     v.IsInput,
			Constraints: convertConstraintsToResponse(v.Constraints),
		}

		if v.Expression != nil {
			item.Expression.SetTo(*v.Expression)
		}

		return item
	})
}

func convertConstraintsToResponse(constraints []version_get_domain.Constraint) []api.TemplateGetByIDVersionVariablesItemConstraintsItem {
	return lo.Map(constraints, func(c version_get_domain.Constraint, _ int) api.TemplateGetByIDVersionVariablesItemConstraintsItem {
		return api.TemplateGetByIDVersionVariablesItemConstraintsItem{
			ID:         c.ID,
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
}
