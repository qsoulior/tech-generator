package template_users_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateUsers(ctx context.Context, params api.TemplateUsersParams) (api.TemplateUsersRes, error) {
	in := domain.TemplateUserListIn{
		UserID:     params.XUserID,
		TemplateID: params.TemplateID,
	}

	users, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("template user list usecase: %w", err)
	}

	resp := api.TemplateUsersResponse{
		Users: lo.Map(users, func(u domain.TemplateUser, _ int) api.TemplateUsersResponseUsersItem {
			return api.TemplateUsersResponseUsersItem{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
				Role:  api.TemplateUsersResponseUsersItemRole(u.Role),
			}
		}),
	}
	return &resp, nil
}
