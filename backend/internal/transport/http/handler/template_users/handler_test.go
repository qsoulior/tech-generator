package template_users_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

func TestHandler_TemplateUsers_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateUsersParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateUserListIn{UserID: 1, TemplateID: 10}).
		Return([]domain.TemplateUser{
			{ID: 2, Name: "alice", Email: "alice@example.com", Role: user_domain.RoleRead},
			{ID: 3, Name: "bob", Email: "bob@example.com", Role: user_domain.RoleWrite},
		}, nil)

	handler := New(usecase)
	got, err := handler.TemplateUsers(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateUsersResponse)
	require.True(t, ok, "expected *api.TemplateUsersResponse, got %T", got)
	require.Len(t, resp.Users, 2)
	require.Equal(t, api.TemplateUsersResponseUsersItemRoleRead, resp.Users[0].Role)
	require.Equal(t, api.TemplateUsersResponseUsersItemRoleWrite, resp.Users[1].Role)
}

func TestHandler_TemplateUsers_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateUsersParams{TemplateID: 10, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrTemplateNotFound},
		{name: "Invalid", err: domain.ErrTemplateInvalid},
		{name: "Wrapped", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.TemplateUsers(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateUsers_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateUsersParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateUsers(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template user list usecase")
	require.ErrorContains(t, err, "boom")
}
