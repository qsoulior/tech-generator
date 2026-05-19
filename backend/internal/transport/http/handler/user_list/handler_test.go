package user_list_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

func TestHandler_UserList_Success(t *testing.T) {
	ctx := context.Background()
	params := api.UserListParams{XUserID: 1, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := domain.UserListIn{
		Page:   1,
		Size:   10,
		Filter: domain.UserListFilter{ExcludeUserID: 1},
	}
	out := domain.UserListOut{
		Users:      []domain.User{{ID: 2, Name: "alice", Email: "alice@example.com"}},
		TotalUsers: 1,
		TotalPages: 1,
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&out, nil)

	handler := New(usecase)
	got, err := handler.UserList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.UserListResponse)
	require.True(t, ok, "expected *api.UserListResponse, got %T", got)
	require.Equal(t, int64(1), resp.TotalUsers)
	require.Equal(t, int64(1), resp.TotalPages)
	require.Len(t, resp.Users, 1)
	require.Equal(t, "alice", resp.Users[0].Name)
}

func TestHandler_UserList_PassesFilter(t *testing.T) {
	ctx := context.Background()
	params := api.UserListParams{XUserID: 1, Page: 2, Size: 5}
	params.UserName.SetTo("foo")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userName := "foo"
	in := domain.UserListIn{
		Page: 2,
		Size: 5,
		Filter: domain.UserListFilter{
			ExcludeUserID: 1,
			UserName:      &userName,
		},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&domain.UserListOut{}, nil)

	handler := New(usecase)
	_, err := handler.UserList(ctx, params)
	require.NoError(t, err)
}

func TestHandler_UserList_ValidationError(t *testing.T) {
	ctx := context.Background()
	params := api.UserListParams{XUserID: 1, Page: 0, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("page", errors.New("invalid"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.UserList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_UserList_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.UserListParams{XUserID: 1, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.UserList(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "user list usecase")
	require.ErrorContains(t, err, "boom")
}
