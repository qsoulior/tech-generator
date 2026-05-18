package template_create_from_default_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/domain"
)

func TestHandler_TemplateCreateFromDefault_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateCreateFromDefaultParams{XUserID: 1}
	req := &api.TemplateCreateFromDefaultRequest{
		SourceTemplateID: 5,
		ProjectID:        2,
		Name:             "copy",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := domain.TemplateCreateFromDefaultIn{
		AuthorID:         1,
		ProjectID:        2,
		SourceTemplateID: 5,
		Name:             "copy",
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&domain.TemplateCreateFromDefaultOut{ID: 42}, nil)

	handler := New(usecase)
	got, err := handler.TemplateCreateFromDefault(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateCreateFromDefaultResponse)
	require.True(t, ok, "expected *api.TemplateCreateFromDefaultResponse, got %T", got)
	require.Equal(t, int64(42), resp.ID)
}

func TestHandler_TemplateCreateFromDefault_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateCreateFromDefaultParams{XUserID: 1}
	req := &api.TemplateCreateFromDefaultRequest{SourceTemplateID: 5, ProjectID: 2, Name: "copy"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, domain.ErrSourceTemplateNotFound)

	handler := New(usecase)
	got, err := handler.TemplateCreateFromDefault(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, domain.ErrSourceTemplateNotFound.Error(), resp.Message)
}

func TestHandler_TemplateCreateFromDefault_ValidationError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateCreateFromDefaultParams{XUserID: 1}
	req := &api.TemplateCreateFromDefaultRequest{SourceTemplateID: 5, ProjectID: 2, Name: ""}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.TemplateCreateFromDefault(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TemplateCreateFromDefault_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateCreateFromDefaultParams{XUserID: 1}
	req := &api.TemplateCreateFromDefaultRequest{SourceTemplateID: 5, ProjectID: 2, Name: "copy"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateCreateFromDefault(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template create from default usecase")
	require.ErrorContains(t, err, "boom")
}
