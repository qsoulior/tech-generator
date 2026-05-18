package template_get_by_id_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

func TestHandler_TemplateGetByID_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	expr := "x+1"
	out := &domain.TemplateGetByIDOut{
		Name: "tmpl",
		Version: &version_get_domain.Version{
			ID:        5,
			Number:    2,
			CreatedAt: createdAt,
			Data:      []byte("data"),
			Variables: []version_get_domain.Variable{{
				ID:         11,
				Name:       "v1",
				Type:       variable_domain.TypeString,
				Expression: &expr,
				IsInput:    true,
				Constraints: []version_get_domain.Constraint{{
					ID:         21,
					VariableID: 11,
					Name:       "c1",
					Expression: "len(x)>0",
					IsActive:   true,
				}},
			}},
		},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1}).
		Return(out, nil)

	handler := New(usecase)
	got, err := handler.TemplateGetByID(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateGetByIDResponse)
	require.True(t, ok, "expected *api.TemplateGetByIDResponse, got %T", got)
	require.Equal(t, "tmpl", resp.Name)

	version, ok := resp.Version.Get()
	require.True(t, ok)
	require.Equal(t, int64(5), version.ID)
	require.Equal(t, int64(2), version.Number)
	require.Equal(t, createdAt, version.CreatedAt)
	require.Equal(t, []byte("data"), version.Data)
	require.Len(t, version.Variables, 1)
	require.Equal(t, int64(11), version.Variables[0].ID)
	require.Equal(t, "v1", version.Variables[0].Name)
	gotExpr, ok := version.Variables[0].Expression.Get()
	require.True(t, ok)
	require.Equal(t, expr, gotExpr)
	require.Len(t, version.Variables[0].Constraints, 1)
	require.Equal(t, int64(21), version.Variables[0].Constraints[0].ID)
}

func TestHandler_TemplateGetByID_SuccessNoVersion(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1}).
		Return(&domain.TemplateGetByIDOut{Name: "tmpl"}, nil)

	handler := New(usecase)
	got, err := handler.TemplateGetByID(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateGetByIDResponse)
	require.True(t, ok)
	require.Equal(t, "tmpl", resp.Name)
	require.False(t, resp.Version.IsSet())
}

func TestHandler_TemplateGetByID_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetByIDParams{TemplateID: 10, XUserID: 1}

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
			got, err := handler.TemplateGetByID(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateGetByID_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateGetByIDParams{TemplateID: 10, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateGetByID(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template get by id usecase")
	require.ErrorContains(t, err, "boom")
}
