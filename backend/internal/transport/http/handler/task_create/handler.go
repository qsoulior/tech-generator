package task_create_handler

import (
	"context"
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TaskCreate(ctx context.Context, req *api.TaskCreateRequest, params api.TaskCreateParams) (api.TaskCreateRes, error) {
	in := domain.TaskCreateIn{
		VersionID: req.VersionID,
		CreatorID: params.XUserID,
		Payload:   req.Payload,
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("task create usecase: %w", err)
	}

	return &api.TaskCreateCreated{}, nil
}
