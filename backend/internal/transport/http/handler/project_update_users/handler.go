package project_update_users_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectUpdateUsers(ctx context.Context, req *api.ProjectUpdateUsersRequest, params api.ProjectUpdateUsersParams) (api.ProjectUpdateUsersRes, error) {
	in := domain.ProjectUserUpdateIn{
		UserID:    params.XUserID,
		ProjectID: params.ProjectID,
		Users: lo.Map(req.Users, func(u api.ProjectUpdateUsersRequestUsersItem, _ int) domain.ProjectUser {
			return domain.ProjectUser{
				ID:   u.ID,
				Role: user_domain.Role(u.Role),
			}
		}),
	}

	err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("project user update usecase: %w", err)
	}

	return &api.ProjectUpdateUsersNoContent{}, nil
}
