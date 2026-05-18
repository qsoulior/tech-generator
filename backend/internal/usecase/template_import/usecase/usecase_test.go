package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	t.Run("IsAuthor/NoVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		in := domain.TemplateImportIn{
			AuthorID:  1,
			ProjectID: 2,
			Name:      "test",
			Version:   nil,
		}

		projectRepo := NewMockprojectRepository(ctrl)
		templateRepo := NewMocktemplateRepository(ctrl)
		versionCreateService := NewMockversionCreateService(ctrl)

		projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&domain.Project{AuthorID: 1}, nil)
		templateRepo.EXPECT().Create(ctx, domain.Template{
			Name:      "test",
			IsDefault: false,
			ProjectID: 2,
			AuthorID:  1,
		}).Return(int64(42), nil)

		usecase := New(projectRepo, templateRepo, versionCreateService)
		got, err := usecase.Handle(ctx, in)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, int64(42), got.ID)
	})

	t.Run("IsWriter/WithVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expr := "x + 1"
		in := domain.TemplateImportIn{
			AuthorID:  1,
			ProjectID: 2,
			Name:      "test",
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

		project := domain.Project{
			AuthorID: 3,
			Users:    []domain.ProjectUser{{ID: 1, Role: user_domain.RoleWrite}},
		}

		projectRepo := NewMockprojectRepository(ctrl)
		templateRepo := NewMocktemplateRepository(ctrl)
		versionCreateService := NewMockversionCreateService(ctrl)

		projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&project, nil)
		templateRepo.EXPECT().Create(ctx, domain.Template{
			Name:      "test",
			IsDefault: false,
			ProjectID: 2,
			AuthorID:  1,
		}).Return(int64(42), nil)

		versionCreateService.EXPECT().Handle(ctx, version_create_domain.VersionCreateIn{
			AuthorID:   1,
			TemplateID: 42,
			Data:       []byte("body"),
			Variables: []version_create_domain.Variable{
				{
					Name:       "x",
					Type:       variable_domain.TypeInteger,
					Expression: &expr,
					IsInput:    false,
					Constraints: []version_create_domain.Constraint{
						{Name: "positive", Expression: "x > 0", IsActive: true},
					},
				},
			},
		}).Return(int64(100), nil)

		usecase := New(projectRepo, templateRepo, versionCreateService)
		got, err := usecase.Handle(ctx, in)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, int64(42), got.ID)
	})
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateImportIn
		setup func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService)
		want  string
	}{
		{
			name:  "in_Validate/Name",
			in:    domain.TemplateImportIn{Name: ""},
			setup: func(*MockprojectRepository, *MocktemplateRepository, *MockversionCreateService) {},
			want:  domain.ErrValueEmpty.Error(),
		},
		{
			name: "in_Validate/VariableType",
			in: domain.TemplateImportIn{
				Name:      "test",
				ProjectID: 2,
				AuthorID:  1,
				Version: &domain.Version{
					Variables: []domain.Variable{{Type: "bogus"}},
				},
			},
			setup: func(*MockprojectRepository, *MocktemplateRepository, *MockversionCreateService) {},
			want:  domain.ErrValueInvalid.Error(),
		},
		{
			name: "projectRepo_GetByID",
			in:   domain.TemplateImportIn{Name: "test", ProjectID: 2, AuthorID: 1},
			setup: func(projectRepo *MockprojectRepository, _ *MocktemplateRepository, _ *MockversionCreateService) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in:   domain.TemplateImportIn{Name: "test", ProjectID: 2, AuthorID: 1},
			setup: func(projectRepo *MockprojectRepository, _ *MocktemplateRepository, _ *MockversionCreateService) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in:   domain.TemplateImportIn{Name: "test", ProjectID: 2, AuthorID: 1},
			setup: func(projectRepo *MockprojectRepository, _ *MocktemplateRepository, _ *MockversionCreateService) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: gofakeit.Int64()}, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "templateRepo_Create",
			in:   domain.TemplateImportIn{Name: "test", ProjectID: 2, AuthorID: 1},
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository, _ *MockversionCreateService) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				templateRepo.EXPECT().Create(ctx, gomock.Any()).Return(int64(0), errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "versionCreateService_Handle",
			in: domain.TemplateImportIn{
				Name:      "test",
				ProjectID: 2,
				AuthorID:  1,
				Version: &domain.Version{
					Data: []byte("body"),
					Variables: []domain.Variable{
						{Name: "x", Type: variable_domain.TypeString, IsInput: true},
					},
				},
			},
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				templateRepo.EXPECT().Create(ctx, gomock.Any()).Return(int64(42), nil)
				versionCreateService.EXPECT().Handle(ctx, gomock.Any()).Return(int64(0), errors.New("test3"))
			},
			want: "test3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			templateRepo := NewMocktemplateRepository(ctrl)
			versionCreateService := NewMockversionCreateService(ctrl)
			tt.setup(projectRepo, templateRepo, versionCreateService)

			usecase := New(projectRepo, templateRepo, versionCreateService)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}

