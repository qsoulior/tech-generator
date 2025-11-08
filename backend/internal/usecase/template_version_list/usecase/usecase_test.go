package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	createdAt := gofakeit.Date()

	tests := []struct {
		name  string
		in    domain.TemplateVersionListIn
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository)
		want  domain.TemplateVersionListOut
	}{
		{
			name: "IsAuthor",
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				versions := []domain.TemplateVersion{
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
				templateVersionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(versions, nil)
			},
			want: domain.TemplateVersionListOut{
				Versions: []domain.TemplateVersion{
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
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				versions := []domain.TemplateVersion{}
				templateVersionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(versions, nil)
			},
			want: domain.TemplateVersionListOut{
				Versions: []domain.TemplateVersion{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateVersionRepo := NewMocktemplateVersionRepository(ctrl)
			tt.setup(templateRepo, templateVersionRepo)

			usecase := New(templateRepo, templateVersionRepo)

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
		in    domain.TemplateVersionListIn
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateVersionRepo_ListByTemplateID",
			in:   domain.TemplateVersionListIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().ListByTemplateID(ctx, int64(10)).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateVersionRepo := NewMocktemplateVersionRepository(ctrl)
			tt.setup(templateRepo, templateVersionRepo)

			usecase := New(templateRepo, templateVersionRepo)

			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
