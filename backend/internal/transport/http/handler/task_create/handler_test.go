package task_create_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

func TestHandler_TaskCreate_Success(t *testing.T) {
	ctx := context.Background()
	payload := api.TaskCreateRequestPayload{"key": "value"}
	req := &api.TaskCreateRequest{VersionID: 7, Payload: payload}
	params := api.TaskCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TaskCreateIn{VersionID: 7, CreatorID: 1, Payload: payload}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.TaskCreate(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.TaskCreateCreated)
	require.True(t, ok, "expected *api.TaskCreateCreated, got %T", got)
}

func TestHandler_TaskCreate_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.TaskCreateRequest{VersionID: 7, Payload: api.TaskCreateRequestPayload{}}
	params := api.TaskCreateParams{XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "VersionNotFound", err: domain.ErrVersionNotFound},
		{name: "VersionInvalid", err: domain.ErrVersionInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(tt.err)

			handler := New(usecase)
			got, err := handler.TaskCreate(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TaskCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.TaskCreateRequest{VersionID: 7, Payload: api.TaskCreateRequestPayload{}}
	params := api.TaskCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TaskCreate(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "task create usecase")
	require.ErrorContains(t, err, "boom")
}
