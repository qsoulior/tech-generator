package template_update_users_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

func TestHandler_TemplateUpdateUsers_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateUsersRequest{
		Users: []api.TemplateUpdateUsersRequestUsersItem{
			{ID: 2, Role: api.TemplateUpdateUsersRequestUsersItemRoleRead},
			{ID: 3, Role: api.TemplateUpdateUsersRequestUsersItemRoleWrite},
		},
	}
	params := api.TemplateUpdateUsersParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateUserUpdateIn{
			UserID:     1,
			TemplateID: 10,
			Users: []domain.TemplateUser{
				{ID: 2, Role: user_domain.RoleRead},
				{ID: 3, Role: user_domain.RoleWrite},
			},
		}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.TemplateUpdateUsers(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.TemplateUpdateUsersNoContent)
	require.True(t, ok, "expected *api.TemplateUpdateUsersNoContent, got %T", got)
}

func TestHandler_TemplateUpdateUsers_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateUsersRequest{}
	params := api.TemplateUpdateUsersParams{TemplateID: 10, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrTemplateNotFound},
		{name: "Invalid", err: domain.ErrTemplateInvalid},
		{name: "RoleInvalid", err: domain.ErrRoleInvalid},
		{name: "Wrapped", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(tt.err)

			handler := New(usecase)
			got, err := handler.TemplateUpdateUsers(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateUpdateUsers_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateUsersRequest{}
	params := api.TemplateUpdateUsersParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateUpdateUsers(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template user update usecase")
	require.ErrorContains(t, err, "boom")
}
