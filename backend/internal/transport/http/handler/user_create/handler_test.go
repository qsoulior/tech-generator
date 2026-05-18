package user_create_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

func TestHandler_UserCreate_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{Name: "alice", Email: "a@b.c", Password: "p"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.UserCreateIn{Name: "alice", Email: "a@b.c", Password: domain.Password("p")}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.UserCreate(ctx, req)
	require.NoError(t, err)

	_, ok := got.(*api.UserCreateCreated)
	require.True(t, ok, "expected *api.UserCreateCreated, got %T", got)
}

func TestHandler_UserCreate_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{Name: "alice", Email: "a@b.c", Password: "p"}

	tests := []struct {
		name string
		err  error
	}{
		{name: "UserExists", err: domain.ErrUserExists},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(tt.err)

			handler := New(usecase)
			got, err := handler.UserCreate(ctx, req)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_UserCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{Name: "", Email: "a@b.c", Password: "p"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(validationErr)

	handler := New(usecase)
	got, err := handler.UserCreate(ctx, req)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_UserCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{Name: "alice", Email: "a@b.c", Password: "p"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.UserCreate(ctx, req)
	require.Nil(t, got)
	require.ErrorContains(t, err, "user create usecase")
	require.ErrorContains(t, err, "boom")
}
