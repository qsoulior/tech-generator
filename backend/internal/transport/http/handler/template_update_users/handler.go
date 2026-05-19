package template_update_users_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TemplateUpdateUsers(ctx context.Context, req *api.TemplateUpdateUsersRequest, params api.TemplateUpdateUsersParams) (api.TemplateUpdateUsersRes, error) {
	in := domain.TemplateUserUpdateIn{
		UserID:     params.XUserID,
		TemplateID: params.TemplateID,
		Users: lo.Map(req.Users, func(u api.TemplateUpdateUsersRequestUsersItem, _ int) domain.TemplateUser {
			return domain.TemplateUser{
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
		return nil, fmt.Errorf("template user update usecase: %w", err)
	}

	return &api.TemplateUpdateUsersNoContent{}, nil
}
