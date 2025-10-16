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
		setup func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.TemplateCreateIn{
				Name:     "test",
				FolderID: 2,
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folder := domain.Folder{
					AuthorID:     1,
					RootAuthorID: 3,
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folder, nil)

				template := domain.Template{
					Name:         "test",
					IsDefault:    false,
					FolderID:     2,
					AuthorID:     1,
					RootAuthorID: 3,
				}
				templateRepo.EXPECT().Create(ctx, template).Return(nil)
			},
		},
		{
			name: "IsRootAuthor",
			in: domain.TemplateCreateIn{
				Name:     "test",
				FolderID: 2,
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folder := domain.Folder{
					AuthorID:     3,
					RootAuthorID: 1,
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folder, nil)

				template := domain.Template{
					Name:         "test",
					IsDefault:    false,
					FolderID:     2,
					AuthorID:     1,
					RootAuthorID: 1,
				}
				templateRepo.EXPECT().Create(ctx, template).Return(nil)
			},
		},
		{
			name: "IsWriter",
			in: domain.TemplateCreateIn{
				Name:     "test",
				FolderID: 2,
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folder := domain.Folder{
					AuthorID:     3,
					RootAuthorID: 4,
					Users:        []domain.FolderUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folder, nil)

				template := domain.Template{
					Name:         "test",
					IsDefault:    false,
					FolderID:     2,
					AuthorID:     1,
					RootAuthorID: 4,
				}
				templateRepo.EXPECT().Create(ctx, template).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(folderRepo, templateRepo)

			usecase := New(folderRepo, templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository)
		in    domain.TemplateCreateIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {},
			in:    domain.TemplateCreateIn{Name: ""},
			want:  domain.ErrValueEmpty.Error(),
		},
		{
			name: "folderRepo_GetByID",
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   domain.TemplateCreateIn{Name: "test", FolderID: 2, AuthorID: 1},
			want: "test1",
		},
		{
			name: "domain_ErrFolderNotFound",
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			in:   domain.TemplateCreateIn{Name: "test", FolderID: 2, AuthorID: 1},
			want: domain.ErrFolderNotFound.Error(),
		},
		{
			name: "domain_ErrFolderInvalid",
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folder := domain.Folder{
					AuthorID:     gofakeit.Int64(),
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
			},
			in:   domain.TemplateCreateIn{Name: "test", FolderID: 2, AuthorID: 1},
			want: domain.ErrFolderInvalid.Error(),
		},
		{
			name: "templateRepo_Create",
			setup: func(folderRepo *MockfolderRepository, templateRepo *MocktemplateRepository) {
				folder := domain.Folder{
					AuthorID:     1,
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
				templateRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test2"))
			},
			in:   domain.TemplateCreateIn{Name: "test", FolderID: 2, AuthorID: 1},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(folderRepo, templateRepo)

			usecase := New(folderRepo, templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
