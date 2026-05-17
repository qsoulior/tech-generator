package project_update_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/domain"
)

func TestHandler_ProjectUpdateByID_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateRequest{Name: "new"}
	params := api.ProjectUpdateByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.ProjectUpdateByID(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.ProjectUpdateByIDNoContent)
	require.True(t, ok, "expected *api.ProjectUpdateByIDNoContent, got %T", got)
}

func TestHandler_ProjectUpdateByID_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateRequest{Name: "new"}
	params := api.ProjectUpdateByIDParams{ProjectID: 10, XUserID: 1}

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
				Handle(ctx, domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"}).
				Return(tt.err)

			handler := New(usecase)
			got, err := handler.ProjectUpdateByID(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_ProjectUpdateByID_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateRequest{Name: ""}
	params := api.ProjectUpdateByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: ""}).
		Return(validationErr)

	handler := New(usecase)
	got, err := handler.ProjectUpdateByID(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_ProjectUpdateByID_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectUpdateRequest{Name: "new"}
	params := api.ProjectUpdateByIDParams{ProjectID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"}).
		Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectUpdateByID(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project update by id usecase")
	require.ErrorContains(t, err, "boom")
}
