package version_list_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

func TestHandler_VersionList_Success(t *testing.T) {
	ctx := context.Background()
	params := api.VersionListParams{XUserID: 1, TemplateID: 3}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	out := &domain.VersionListOut{
		Versions: []domain.Version{{
			ID:         5,
			Number:     2,
			AuthorName: "alice",
			CreatedAt:  createdAt,
		}},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.VersionListIn{TemplateID: 3, UserID: 1}).
		Return(out, nil)

	handler := New(usecase)
	got, err := handler.VersionList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.VersionListResponse)
	require.True(t, ok, "expected *api.VersionListResponse, got %T", got)
	require.Len(t, resp.Versions, 1)
	require.Equal(t, int64(5), resp.Versions[0].ID)
	require.Equal(t, int64(2), resp.Versions[0].Number)
	require.Equal(t, "alice", resp.Versions[0].AuthorName)
	require.Equal(t, createdAt, resp.Versions[0].CreatedAt)
}

func TestHandler_VersionList_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.VersionListParams{XUserID: 1, TemplateID: 3}

	tests := []struct {
		name string
		err  error
	}{
		{name: "TemplateNotFound", err: domain.ErrTemplateNotFound},
		{name: "TemplateInvalid", err: domain.ErrTemplateInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.VersionList(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_VersionList_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.VersionListParams{XUserID: 1, TemplateID: 3}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.VersionList(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "version list usecase")
	require.ErrorContains(t, err, "boom")
}
