package template_get_meta_by_id_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_meta_by_id/domain"
)

func TestHandler_TemplateGetMetaByID_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetMetaByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateGetMetaByIDIn{TemplateID: 10, UserID: 1}).
		Return(&domain.TemplateGetMetaByIDOut{Name: "tmpl"}, nil)

	handler := New(usecase)
	got, err := handler.TemplateGetMetaByID(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateGetMetaByIDResponse)
	require.True(t, ok, "expected *api.TemplateGetMetaByIDResponse, got %T", got)
	require.Equal(t, "tmpl", resp.Name)
}

func TestHandler_TemplateGetMetaByID_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetMetaByIDParams{TemplateID: 10, XUserID: 1}

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
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.TemplateGetMetaByID(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateGetMetaByID_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetMetaByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateGetMetaByID(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template get meta by id usecase")
	require.ErrorContains(t, err, "boom")
}
