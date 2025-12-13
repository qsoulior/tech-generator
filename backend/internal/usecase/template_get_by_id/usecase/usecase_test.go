package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	in := domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService)
		want  domain.TemplateGetByIDOut
	}{
		{
			name: "IsProjectAuthor",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					Name:            "test",
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(&version_get_domain.Version{ID: 20}, nil)
			},
			want: domain.TemplateGetByIDOut{Name: "test", Version: &version_get_domain.Version{ID: 20}},
		},
		{
			name: "IsAuthor/NoLastVersion",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					Name:            "test",
					LastVersionID:   nil,
					AuthorID:        1,
					ProjectAuthorID: 2,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetByIDOut{Name: "test", Version: nil},
		},
		{
			name: "IsReader/NoLastVersion",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					Name:            "test",
					LastVersionID:   nil,
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleRead}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetByIDOut{Name: "test", Version: nil},
		},
		{
			name: "IsWriter/NoLastVersion",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					Name:            "test",
					LastVersionID:   nil,
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetByIDOut{Name: "test", Version: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)

			tt.setup(templateRepo, versionGetService)

			usecase := New(templateRepo, versionGetService)

			got, err := usecase.Handle(ctx, in)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	in := domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 3,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "versionGetService_Handle",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)

			tt.setup(templateRepo, versionGetService)

			usecase := New(templateRepo, versionGetService)

			_, err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
