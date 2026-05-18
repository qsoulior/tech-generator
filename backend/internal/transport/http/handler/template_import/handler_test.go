package template_import_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

func TestHandler_TemplateImport_Success(t *testing.T) {
	ctx := context.Background()
	expr := "x + 1"

	req := &api.TemplateImportRequest{
		ProjectID: 3,
		Template: api.TemplateImportPayload{
			Name: "tmpl",
			Version: api.NewOptTemplateImportVersion(api.TemplateImportVersion{
				Data: []byte("body"),
				Variables: []api.TemplateImportVersionVariablesItem{
					{
						Name:       "x",
						Type:       api.TemplateImportVersionVariablesItemTypeInteger,
						Expression: api.NewOptString(expr),
						IsInput:    false,
						Constraints: []api.TemplateImportVersionVariablesItemConstraintsItem{
							{Name: "positive", Expression: "x > 0", IsActive: true},
						},
					},
				},
			}),
		},
	}
	params := api.TemplateImportParams{XUserID: 1}

	want := domain.TemplateImportIn{
		AuthorID:  1,
		ProjectID: 3,
		Name:      "tmpl",
		Version: &domain.Version{
			Data: []byte("body"),
			Variables: []domain.Variable{
				{
					Name:       "x",
					Type:       variable_domain.TypeInteger,
					Expression: &expr,
					IsInput:    false,
					Constraints: []domain.Constraint{
						{Name: "positive", Expression: "x > 0", IsActive: true},
					},
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, want).Return(&domain.TemplateImportOut{ID: 42}, nil)

	handler := New(usecase)
	got, err := handler.TemplateImport(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateImportResponse)
	require.True(t, ok, "expected *api.TemplateImportResponse, got %T", got)
	require.Equal(t, int64(42), resp.ID)
}

func TestHandler_TemplateImport_NoVersion(t *testing.T) {
	ctx := context.Background()

	req := &api.TemplateImportRequest{
		ProjectID: 3,
		Template:  api.TemplateImportPayload{Name: "tmpl"},
	}
	params := api.TemplateImportParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, domain.TemplateImportIn{
		AuthorID:  1,
		ProjectID: 3,
		Name:      "tmpl",
	}).Return(&domain.TemplateImportOut{ID: 7}, nil)

	handler := New(usecase)
	got, err := handler.TemplateImport(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateImportResponse)
	require.True(t, ok)
	require.Equal(t, int64(7), resp.ID)
}

func TestHandler_TemplateImport_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateImportRequest{ProjectID: 3, Template: api.TemplateImportPayload{Name: "tmpl"}}
	params := api.TemplateImportParams{XUserID: 1}

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
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.TemplateImport(ctx, req, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TemplateImport_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateImportRequest{ProjectID: 3, Template: api.TemplateImportPayload{Name: ""}}
	params := api.TemplateImportParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.TemplateImport(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TemplateImport_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.TemplateImportRequest{ProjectID: 3, Template: api.TemplateImportPayload{Name: "tmpl"}}
	params := api.TemplateImportParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateImport(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template import usecase")
	require.ErrorContains(t, err, "boom")
}
