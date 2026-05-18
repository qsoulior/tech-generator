package version_create_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

func TestHandler_VersionCreate_Success(t *testing.T) {
	ctx := context.Background()
	expr := "x+1"
	req := &api.VersionCreateRequest{
		TemplateID: 3,
		Data:       []byte("data"),
		Variables: []api.VersionCreateRequestVariablesItem{{
			Name:       "v",
			Type:       api.VersionCreateRequestVariablesItemType(variable_domain.TypeString),
			Expression: api.NewOptString(expr),
			IsInput:    true,
			Constraints: []api.VersionCreateRequestVariablesItemConstraintsItem{{
				Name:       "c",
				Expression: "len(x)>0",
				IsActive:   true,
			}},
		}},
	}
	params := api.VersionCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := version_create_domain.VersionCreateIn{
		AuthorID:   1,
		TemplateID: 3,
		Data:       []byte("data"),
		Variables: []version_create_domain.Variable{{
			Name:       "v",
			Type:       variable_domain.TypeString,
			Expression: &expr,
			IsInput:    true,
			Constraints: []version_create_domain.Constraint{{
				Name:       "c",
				Expression: "len(x)>0",
				IsActive:   true,
			}},
		}},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(int64(42), nil)

	handler := New(usecase)
	got, err := handler.VersionCreate(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.VersionCreateResponse)
	require.True(t, ok, "expected *api.VersionCreateResponse, got %T", got)
	require.Equal(t, int64(42), resp.ID)
}

func TestHandler_VersionCreate_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateRequest{TemplateID: 3, Data: []byte("d")}
	params := api.VersionCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseErr := error_domain.NewBaseError("custom")
	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(int64(0), baseErr)

	handler := New(usecase)
	got, err := handler.VersionCreate(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, baseErr.Error(), resp.Message)
}

func TestHandler_VersionCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateRequest{TemplateID: 3, Data: []byte("d")}
	params := api.VersionCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("variables.0.type", errors.New("invalid"))
	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(int64(0), validationErr)

	handler := New(usecase)
	got, err := handler.VersionCreate(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_VersionCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.VersionCreateRequest{TemplateID: 3, Data: []byte("d")}
	params := api.VersionCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(int64(0), errors.New("boom"))

	handler := New(usecase)
	got, err := handler.VersionCreate(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "version create usecase")
	require.ErrorContains(t, err, "boom")
}
