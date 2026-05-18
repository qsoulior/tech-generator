package template_create_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

func TestHandler_TemplateCreate_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateCreateRequest{Name: "tmpl", ProjectID: 3}
	params := api.TemplateCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateCreateIn{Name: "tmpl", ProjectID: 3, AuthorID: 1}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.TemplateCreate(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.TemplateCreateCreated)
	require.True(t, ok, "expected *api.TemplateCreateCreated, got %T", got)
}

func TestHandler_TemplateCreate_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateCreateRequest{Name: "tmpl", ProjectID: 3}
	params := api.TemplateCreateParams{XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "ProjectNotFound", err: domain.ErrProjectNotFound},
		{name: "ProjectInvalid", err: domain.ErrProjectInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(tt.err)

			handler := New(usecase)
			got, err := handler.TemplateCreate(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateCreateRequest{Name: "", ProjectID: 3}
	params := api.TemplateCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(validationErr)

	handler := New(usecase)
	got, err := handler.TemplateCreate(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TemplateCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateCreateRequest{Name: "tmpl", ProjectID: 3}
	params := api.TemplateCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateCreate(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template create usecase")
	require.ErrorContains(t, err, "boom")
}
