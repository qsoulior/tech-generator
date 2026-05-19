package user_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) UserList(ctx context.Context, params api.UserListParams) (api.UserListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}
		return nil, fmt.Errorf("user list usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.UserListParams) domain.UserListIn {
	filter := domain.UserListFilter{
		ExcludeUserID: params.XUserID,
	}

	if userName, ok := params.UserName.Get(); ok {
		filter.UserName = &userName
	}

	return domain.UserListIn{
		Page:   params.Page,
		Size:   params.Size,
		Filter: filter,
	}
}

func convertOutToResponse(out domain.UserListOut) api.UserListResponse {
	return api.UserListResponse{
		Users: lo.Map(out.Users, func(u domain.User, _ int) api.UserListResponseUsersItem {
			return api.UserListResponseUsersItem{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
			}
		}),
		TotalUsers: out.TotalUsers,
		TotalPages: out.TotalPages,
	}
}
