package version_create_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) VersionCreate(ctx context.Context, req *api.VersionCreateRequest, params api.VersionCreateParams) (api.VersionCreateRes, error) {
	err := h.usecase.Handle(ctx, convertRequestToIn(req, params))
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("version create usecase: %w", err)
	}

	return &api.VersionCreateCreated{}, nil
}

func convertRequestToIn(req *api.VersionCreateRequest, params api.VersionCreateParams) version_create_domain.VersionCreateIn {
	return version_create_domain.VersionCreateIn{
		AuthorID:   params.XUserID,
		TemplateID: req.TemplateID,
		Data:       req.Data,
		Variables:  convertVariablesToIn(req.Variables),
	}
}

func convertVariablesToIn(variables []api.VersionCreateRequestVariablesItem) []version_create_domain.Variable {
	return lo.Map(variables, func(v api.VersionCreateRequestVariablesItem, _ int) version_create_domain.Variable {
		variable := version_create_domain.Variable{
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

func convertConstraintsToIn(constraints []api.VersionCreateRequestVariablesItemConstraintsItem) []version_create_domain.Constraint {
	return lo.Map(constraints, func(c api.VersionCreateRequestVariablesItemConstraintsItem, _ int) version_create_domain.Constraint {
		return version_create_domain.Constraint{
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
}
