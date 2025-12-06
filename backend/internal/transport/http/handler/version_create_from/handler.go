package version_create_from_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) VersionCreateFrom(ctx context.Context, req *api.VersionCreateFromRequest, params api.VersionCreateFromParams) (api.VersionCreateFromRes, error) {
	err := h.usecase.Handle(ctx, convertRequestToIn(req, params))
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("version create from usecase: %w", err)
	}

	return &api.VersionCreateFromCreated{}, nil
}

func convertRequestToIn(req *api.VersionCreateFromRequest, params api.VersionCreateFromParams) domain.VersionCreateFromIn {
	return domain.VersionCreateFromIn{
		AuthorID:   params.XUserID,
		TemplateID: req.TemplateID,
		VersionID:  req.VersionID,
	}
}
