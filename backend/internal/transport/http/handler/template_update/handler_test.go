package template_update_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_update/domain"
)

func TestHandler_TemplateUpdateByID_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateRequest{Name: "new"}
	params := api.TemplateUpdateByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.TemplateUpdateByID(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.TemplateUpdateByIDNoContent)
	require.True(t, ok, "expected *api.TemplateUpdateByIDNoContent, got %T", got)
}

func TestHandler_TemplateUpdateByID_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateRequest{Name: "new"}
	params := api.TemplateUpdateByIDParams{TemplateID: 10, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrTemplateNotFound},
		{name: "Invalid", err: domain.ErrTemplateInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().
				Handle(ctx, domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"}).
				Return(tt.err)

			handler := New(usecase)
			got, err := handler.TemplateUpdateByID(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateUpdateByID_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateRequest{Name: ""}
	params := api.TemplateUpdateByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: ""}).
		Return(validationErr)

	handler := New(usecase)
	got, err := handler.TemplateUpdateByID(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TemplateUpdateByID_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateUpdateRequest{Name: "new"}
	params := api.TemplateUpdateByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"}).
		Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateUpdateByID(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template update by id usecase")
	require.ErrorContains(t, err, "boom")
}
