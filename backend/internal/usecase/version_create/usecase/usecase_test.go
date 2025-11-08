package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	in := version_create_domain.VersionCreateIn{
		AuthorID:   1,
		TemplateID: 10,
		Data:       []byte{1, 2, 3},
		Variables:  []version_create_domain.Variable{},
	}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService)
	}{
		{
			name: "IsAuthor",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionCreateService.EXPECT().Handle(ctx, in).Return(nil)
			},
		},
		{
			name: "IsRootAuthor",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionCreateService.EXPECT().Handle(ctx, in).Return(nil)
			},
		},
		{
			name: "IsWriter",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				template := domain.Template{
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionCreateService.EXPECT().Handle(ctx, in).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionCreateService := NewMockversionCreateService(ctrl)

			tt.setup(templateRepo, versionCreateService)

			usecase := New(templateRepo, versionCreateService)
			err := usecase.Handle(ctx, in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	in := version_create_domain.VersionCreateIn{
		AuthorID:   1,
		TemplateID: 10,
		Data:       []byte{1, 2, 3},
		Variables:  []version_create_domain.Variable{},
	}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository, versionCreateService *MockversionCreateService) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionCreateService.EXPECT().Handle(ctx, in).Return(errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionCreateService := NewMockversionCreateService(ctrl)

			tt.setup(templateRepo, versionCreateService)

			usecase := New(templateRepo, versionCreateService)
			err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
