package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	createdAt := gofakeit.Date()

	tests := []struct {
		name  string
		in    domain.VersionListIn
		setup func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository)
		want  domain.VersionListOut
	}{
		{
			name: "IsAuthor",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				versions := []domain.Version{
					{
						ID:         21,
						Number:     1,
						AuthorName: "test_1",
						CreatedAt:  createdAt,
					},
					{
						ID:         22,
						Number:     2,
						AuthorName: "test2",
						CreatedAt:  createdAt,
					},
				}
				versionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(versions, nil)
			},
			want: domain.VersionListOut{
				Versions: []domain.Version{
					{
						ID:         21,
						Number:     1,
						AuthorName: "test_1",
						CreatedAt:  createdAt,
					},
					{
						ID:         22,
						Number:     2,
						AuthorName: "test2",
						CreatedAt:  createdAt,
					},
				},
			},
		},
		{
			name: "IsProjectAuthor",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				versions := []domain.Version{}
				versionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(versions, nil)
			},
			want: domain.VersionListOut{
				Versions: []domain.Version{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			tt.setup(templateRepo, versionRepo)

			usecase := New(templateRepo, versionRepo)

			got, err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.VersionListIn
		setup func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "versionRepo_ListByTemplateID",
			in:   domain.VersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				versionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			tt.setup(templateRepo, versionRepo)

			usecase := New(templateRepo, versionRepo)

			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
