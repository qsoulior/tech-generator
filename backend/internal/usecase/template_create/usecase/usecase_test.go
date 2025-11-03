package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateCreateIn
		setup func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.TemplateCreateIn{
				Name:      "test",
				ProjectID: 2,
				AuthorID:  1,
			},
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&project, nil)

				template := domain.Template{
					Name:      "test",
					IsDefault: false,
					ProjectID: 2,
					AuthorID:  1,
				}
				templateRepo.EXPECT().Create(ctx, template).Return(nil)
			},
		},
		{
			name: "IsWriter",
			in: domain.TemplateCreateIn{
				Name:      "test",
				ProjectID: 2,
				AuthorID:  1,
			},
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				project := domain.Project{
					AuthorID: 3,
					Users:    []domain.ProjectUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&project, nil)

				template := domain.Template{
					Name:      "test",
					IsDefault: false,
					ProjectID: 2,
					AuthorID:  1,
				}
				templateRepo.EXPECT().Create(ctx, template).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(projectRepo, templateRepo)

			usecase := New(projectRepo, templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository)
		in    domain.TemplateCreateIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {},
			in:    domain.TemplateCreateIn{Name: ""},
			want:  domain.ErrValueEmpty.Error(),
		},
		{
			name: "projectRepo_GetByID",
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   domain.TemplateCreateIn{Name: "test", ProjectID: 2, AuthorID: 1},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			in:   domain.TemplateCreateIn{Name: "test", ProjectID: 2, AuthorID: 1},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				project := domain.Project{
					AuthorID: gofakeit.Int64(),
					Users:    []domain.ProjectUser{},
				}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
			},
			in:   domain.TemplateCreateIn{Name: "test", ProjectID: 2, AuthorID: 1},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "templateRepo_Create",
			setup: func(projectRepo *MockprojectRepository, templateRepo *MocktemplateRepository) {
				project := domain.Project{
					AuthorID: 1,
					Users:    []domain.ProjectUser{},
				}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
				templateRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test2"))
			},
			in:   domain.TemplateCreateIn{Name: "test", ProjectID: 2, AuthorID: 1},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(projectRepo, templateRepo)

			usecase := New(projectRepo, templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
