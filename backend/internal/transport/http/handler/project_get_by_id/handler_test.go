package project_get_by_id_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

func TestHandler_ProjectGetByID_Success(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectGetByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectGetByIDIn{ProjectID: 10, UserID: 1}).
		Return(&domain.ProjectGetByIDOut{Name: "test", AuthorName: "author"}, nil)

	handler := New(usecase)
	got, err := handler.ProjectGetByID(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.ProjectGetByIDResponse)
	require.True(t, ok, "expected *api.ProjectGetByIDResponse, got %T", got)
	require.Equal(t, "test", resp.Name)
	require.Equal(t, "author", resp.AuthorName)
}

func TestHandler_ProjectGetByID_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectGetByIDParams{ProjectID: 10, XUserID: 1}

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
				Handle(ctx, domain.ProjectGetByIDIn{ProjectID: 10, UserID: 1}).
				Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.ProjectGetByID(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_ProjectGetByID_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectGetByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectGetByIDIn{ProjectID: 10, UserID: 1}).
		Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectGetByID(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project get by id usecase")
	require.ErrorContains(t, err, "boom")
}
