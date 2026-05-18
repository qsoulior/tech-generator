package template_import_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateImport(ctx context.Context, req *api.TemplateImportRequest, params api.TemplateImportParams) (api.TemplateImportRes, error) {
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

		return nil, fmt.Errorf("template import usecase: %w", err)
	}

	return &api.TemplateImportResponse{ID: out.ID}, nil
}

func convertRequestToIn(req *api.TemplateImportRequest, params api.TemplateImportParams) domain.TemplateImportIn {
	in := domain.TemplateImportIn{
		AuthorID:  params.XUserID,
		ProjectID: req.ProjectID,
		Name:      req.Template.Name,
	}

	if version, ok := req.Template.Version.Get(); ok {
		in.Version = &domain.Version{
			Data:      version.Data,
			Variables: convertVariablesToIn(version.Variables),
		}
	}

	return in
}

func convertVariablesToIn(variables []api.TemplateImportVersionVariablesItem) []domain.Variable {
	return lo.Map(variables, func(v api.TemplateImportVersionVariablesItem, _ int) domain.Variable {
		variable := domain.Variable{
			Name:        v.Name,
			Type:        variable_domain.Type(v.Type),
			IsInput:     v.IsInput,
			Constraints: convertConstraintsToIn(v.Constraints),
		}

		if v.Expression.IsSet() {
			variable.Expression = &v.Expression.Value
		}

		return variable
	})
}

func convertConstraintsToIn(constraints []api.TemplateImportVersionVariablesItemConstraintsItem) []domain.Constraint {
	return lo.Map(constraints, func(c api.TemplateImportVersionVariablesItemConstraintsItem, _ int) domain.Constraint {
		return domain.Constraint{
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
}
