package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	in := domain.VersionCreateFromIn{
		AuthorID:   1,
		TemplateID: 10,
		VersionID:  20,
	}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService)
	}{
		{
			name: "IsAuthorID",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				version := version_get_domain.Version{
					TemplateID: 10,
					Data:       []byte{1, 2, 3},
					Variables: []version_get_domain.Variable{
						{
							Name:       "var1",
							Type:       variable_domain.TypeString,
							Expression: lo.ToPtr("expr1"),
							Constraints: []version_get_domain.Constraint{
								{Name: "constraint1", Expression: "expr11", IsActive: false},
							},
						},
						{
							Name:       "var2",
							Type:       variable_domain.TypeInteger,
							Expression: lo.ToPtr("expr2"),
							Constraints: []version_get_domain.Constraint{
								{Name: "constraint2", Expression: "expr12", IsActive: true},
							},
						},
					},
				}
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(&version, nil)

				in := version_create_domain.VersionCreateIn{
					AuthorID:   1,
					TemplateID: 10,
					Data:       []byte{1, 2, 3},
					Variables: []version_create_domain.Variable{
						{
							Name:       "var1",
							Type:       variable_domain.TypeString,
							Expression: lo.ToPtr("expr1"),
							Constraints: []version_create_domain.Constraint{
								{Name: "constraint1", Expression: "expr11", IsActive: false},
							},
						},
						{
							Name:       "var2",
							Type:       variable_domain.TypeInteger,
							Expression: lo.ToPtr("expr2"),
							Constraints: []version_create_domain.Constraint{
								{Name: "constraint2", Expression: "expr12", IsActive: true},
							},
						},
					},
				}
				versionCreateService.EXPECT().Handle(ctx, in).Return(nil)
			},
		},
		{
			name: "IsProjectAuthorID",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				version := version_get_domain.Version{TemplateID: 10}
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(&version, nil)
				versionCreateService.EXPECT().Handle(ctx, gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)
			versionCreateService := NewMockversionCreateService(ctrl)
			tt.setup(templateRepo, versionGetService, versionCreateService)

			usecase := New(templateRepo, versionGetService, versionCreateService)
			err := usecase.Handle(ctx, in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	in := domain.VersionCreateFromIn{
		AuthorID:   1,
		TemplateID: 10,
		VersionID:  20,
	}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "versionGetService_Handle",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "versionGetService_Handle",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				version := version_get_domain.Version{TemplateID: 11}
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(&version, nil)
			},
			want: domain.ErrVersionInvalid.Error(),
		},
		{
			name: "versionCreateService_Handle",
			setup: func(templateRepo *MocktemplateRepository, versionGetService *MockversionGetService, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				version := version_get_domain.Version{TemplateID: 10}
				versionGetService.EXPECT().Handle(ctx, int64(20)).Return(&version, nil)
				versionCreateService.EXPECT().Handle(ctx, gomock.Any()).Return(errors.New("test3"))
			},
			want: "test3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)
			versionCreateService := NewMockversionCreateService(ctrl)
			tt.setup(templateRepo, versionGetService, versionCreateService)

			usecase := New(templateRepo, versionGetService, versionCreateService)
			err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
