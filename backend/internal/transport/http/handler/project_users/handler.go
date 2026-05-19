package project_users_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectUsers(ctx context.Context, params api.ProjectUsersParams) (api.ProjectUsersRes, error) {
	in := domain.ProjectUserListIn{
		UserID:    params.XUserID,
		ProjectID: params.ProjectID,
	}

	users, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("project user list usecase: %w", err)
	}

	resp := api.ProjectUsersResponse{
		Users: lo.Map(users, func(u domain.ProjectUser, _ int) api.ProjectUsersResponseUsersItem {
			return api.ProjectUsersResponseUsersItem{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
				Role:  api.ProjectUsersResponseUsersItemRole(u.Role),
			}
		}),
	}
	return &resp, nil
}
