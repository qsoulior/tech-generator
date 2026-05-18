package version_create_from_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

func TestHandler_VersionCreateFrom_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateFromRequest{TemplateID: 3, VersionID: 5}
	params := api.VersionCreateFromParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.VersionCreateFromIn{AuthorID: 1, TemplateID: 3, VersionID: 5}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.VersionCreateFrom(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.VersionCreateFromCreated)
	require.True(t, ok, "expected *api.VersionCreateFromCreated, got %T", got)
}

func TestHandler_VersionCreateFrom_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateFromRequest{TemplateID: 3, VersionID: 5}
	params := api.VersionCreateFromParams{XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "TemplateNotFound", err: domain.ErrTemplateNotFound},
		{name: "TemplateInvalid", err: domain.ErrTemplateInvalid},
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
			got, err := handler.VersionCreateFrom(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_VersionCreateFrom_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateFromRequest{TemplateID: 3, VersionID: 5}
	params := api.VersionCreateFromParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.VersionCreateFrom(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "version create from usecase")
	require.ErrorContains(t, err, "boom")
}
