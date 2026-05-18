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
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	t.Run("IsAuthor/NoVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		in := domain.TemplateCreateFromDefaultIn{
			AuthorID:         1,
			ProjectID:        2,
			SourceTemplateID: 5,
			Name:             "copy",
		}

		projectRepo := NewMockprojectRepository(ctrl)
		sourceRepo := NewMocksourceTemplateRepository(ctrl)
		newRepo := NewMocknewTemplateRepository(ctrl)
		versionGet := NewMockversionGetService(ctrl)
		versionCreate := NewMockversionCreateService(ctrl)

		projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&domain.Project{AuthorID: 1}, nil)
		sourceRepo.EXPECT().GetByID(ctx, int64(5)).Return(&domain.SourceTemplate{ID: 5, IsDefault: true, LastVersionID: nil}, nil)
		newRepo.EXPECT().Create(ctx, domain.Template{
			Name:      "copy",
			IsDefault: false,
			ProjectID: 2,
			AuthorID:  1,
		}).Return(int64(42), nil)

		usecase := New(projectRepo, sourceRepo, newRepo, versionGet, versionCreate)
		got, err := usecase.Handle(ctx, in)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, int64(42), got.ID)
	})

	t.Run("IsWriter/WithVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expr := "x + 1"
		versionID := int64(100)
		in := domain.TemplateCreateFromDefaultIn{
			AuthorID:         1,
			ProjectID:        2,
			SourceTemplateID: 5,
			Name:             "copy",
		}

		project := domain.Project{
			AuthorID: 3,
			Users:    []domain.ProjectUser{{ID: 1, Role: user_domain.RoleWrite}},
		}
		source := domain.SourceTemplate{ID: 5, IsDefault: true, LastVersionID: &versionID}
		version := version_get_domain.Version{
			ID:         versionID,
			TemplateID: 5,
			Number:     1,
			Data:       []byte("body"),
			Variables: []version_get_domain.Variable{
				{
					ID:         11,
					Name:       "x",
					Type:       variable_domain.TypeInteger,
					Expression: &expr,
					IsInput:    false,
					Constraints: []version_get_domain.Constraint{
						{ID: 21, VariableID: 11, Name: "positive", Expression: "x > 0", IsActive: true},
					},
				},
			},
		}

		projectRepo := NewMockprojectRepository(ctrl)
		sourceRepo := NewMocksourceTemplateRepository(ctrl)
		newRepo := NewMocknewTemplateRepository(ctrl)
		versionGet := NewMockversionGetService(ctrl)
		versionCreate := NewMockversionCreateService(ctrl)

		projectRepo.EXPECT().GetByID(ctx, int64(2)).Return(&project, nil)
		sourceRepo.EXPECT().GetByID(ctx, int64(5)).Return(&source, nil)
		newRepo.EXPECT().Create(ctx, domain.Template{
			Name:      "copy",
			IsDefault: false,
			ProjectID: 2,
			AuthorID:  1,
		}).Return(int64(42), nil)
		versionGet.EXPECT().Handle(ctx, versionID).Return(&version, nil)
		versionCreate.EXPECT().Handle(ctx, version_create_domain.VersionCreateIn{
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
		}).Return(int64(200), nil)

		usecase := New(projectRepo, sourceRepo, newRepo, versionGet, versionCreate)
		got, err := usecase.Handle(ctx, in)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, int64(42), got.ID)
	})
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	validIn := domain.TemplateCreateFromDefaultIn{
		AuthorID:         1,
		ProjectID:        2,
		SourceTemplateID: 5,
		Name:             "copy",
	}

	tests := []struct {
		name  string
		in    domain.TemplateCreateFromDefaultIn
		setup func(
			projectRepo *MockprojectRepository,
			sourceRepo *MocksourceTemplateRepository,
			newRepo *MocknewTemplateRepository,
			versionGet *MockversionGetService,
			versionCreate *MockversionCreateService,
		)
		want string
	}{
		{
			name: "in_Validate/Name",
			in:   domain.TemplateCreateFromDefaultIn{Name: ""},
			setup: func(*MockprojectRepository, *MocksourceTemplateRepository, *MocknewTemplateRepository, *MockversionGetService, *MockversionCreateService) {
			},
			want: domain.ErrValueEmpty.Error(),
		},
		{
			name: "projectRepo_GetByID",
			in:   validIn,
			setup: func(p *MockprojectRepository, _ *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in:   validIn,
			setup: func(p *MockprojectRepository, _ *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in:   validIn,
			setup: func(p *MockprojectRepository, _ *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: gofakeit.Int64()}, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "sourceTemplateRepo_GetByID",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "domain_ErrSourceTemplateNotFound",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrSourceTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrSourceTemplateInvalid",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, _ *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.SourceTemplate{ID: 5, IsDefault: false}, nil)
			},
			want: domain.ErrSourceTemplateInvalid.Error(),
		},
		{
			name: "newTemplateRepo_Create",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, nr *MocknewTemplateRepository, _ *MockversionGetService, _ *MockversionCreateService) {
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.SourceTemplate{ID: 5, IsDefault: true}, nil)
				nr.EXPECT().Create(ctx, gomock.Any()).Return(int64(0), errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "versionGetService_Handle",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, nr *MocknewTemplateRepository, vg *MockversionGetService, _ *MockversionCreateService) {
				versionID := int64(100)
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.SourceTemplate{ID: 5, IsDefault: true, LastVersionID: &versionID}, nil)
				nr.EXPECT().Create(ctx, gomock.Any()).Return(int64(42), nil)
				vg.EXPECT().Handle(ctx, versionID).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "versionCreateService_Handle",
			in:   validIn,
			setup: func(p *MockprojectRepository, sr *MocksourceTemplateRepository, nr *MocknewTemplateRepository, vg *MockversionGetService, vc *MockversionCreateService) {
				versionID := int64(100)
				p.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				sr.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.SourceTemplate{ID: 5, IsDefault: true, LastVersionID: &versionID}, nil)
				nr.EXPECT().Create(ctx, gomock.Any()).Return(int64(42), nil)
				vg.EXPECT().Handle(ctx, versionID).Return(&version_get_domain.Version{Data: []byte("body")}, nil)
				vc.EXPECT().Handle(ctx, gomock.Any()).Return(int64(0), errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			sourceRepo := NewMocksourceTemplateRepository(ctrl)
			newRepo := NewMocknewTemplateRepository(ctrl)
			versionGet := NewMockversionGetService(ctrl)
			versionCreate := NewMockversionCreateService(ctrl)
			tt.setup(projectRepo, sourceRepo, newRepo, versionGet, versionCreate)

			usecase := New(projectRepo, sourceRepo, newRepo, versionGet, versionCreate)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
