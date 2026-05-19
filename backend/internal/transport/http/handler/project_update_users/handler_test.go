package project_update_users_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

func TestHandler_ProjectUpdateUsers_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateUsersRequest{
		Users: []api.ProjectUpdateUsersRequestUsersItem{
			{ID: 2, Role: api.ProjectUpdateUsersRequestUsersItemRoleRead},
			{ID: 3, Role: api.ProjectUpdateUsersRequestUsersItemRoleMaintain},
		},
	}
	params := api.ProjectUpdateUsersParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectUserUpdateIn{
			UserID:    1,
			ProjectID: 10,
			Users: []domain.ProjectUser{
				{ID: 2, Role: user_domain.RoleRead},
				{ID: 3, Role: user_domain.RoleMaintain},
			},
		}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.ProjectUpdateUsers(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.ProjectUpdateUsersNoContent)
	require.True(t, ok, "expected *api.ProjectUpdateUsersNoContent, got %T", got)
}

func TestHandler_ProjectUpdateUsers_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateUsersRequest{}
	params := api.ProjectUpdateUsersParams{ProjectID: 10, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrProjectNotFound},
		{name: "Invalid", err: domain.ErrProjectInvalid},
		{name: "Wrapped", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(tt.err)

			handler := New(usecase)
			got, err := handler.ProjectUpdateUsers(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_ProjectUpdateUsers_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateUsersRequest{}
	params := api.ProjectUpdateUsersParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectUpdateUsers(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project user update usecase")
	require.ErrorContains(t, err, "boom")
}
