package project_delete_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"
)

func TestHandler_ProjectDeleteByID_Success(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectDeleteByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectDeleteIn{ProjectID: 10, UserID: 1}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.ProjectDeleteByID(ctx, params)
	require.NoError(t, err)

	_, ok := got.(*api.ProjectDeleteByIDNoContent)
	require.True(t, ok, "expected *api.ProjectDeleteByIDNoContent, got %T", got)
}

func TestHandler_ProjectDeleteByID_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectDeleteByIDParams{ProjectID: 10, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrProjectNotFound},
		{name: "Invalid", err: domain.ErrProjectInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().
				Handle(ctx, domain.ProjectDeleteIn{ProjectID: 10, UserID: 1}).
				Return(tt.err)

			handler := New(usecase)
			got, err := handler.ProjectDeleteByID(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_ProjectDeleteByID_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectDeleteByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectDeleteIn{ProjectID: 10, UserID: 1}).
		Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectDeleteByID(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project delete by id usecase")
	require.ErrorContains(t, err, "boom")
}
