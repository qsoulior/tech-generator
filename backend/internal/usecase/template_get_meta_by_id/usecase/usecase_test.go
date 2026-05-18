package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_meta_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	in := domain.TemplateGetMetaByIDIn{TemplateID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository)
		want  domain.TemplateGetMetaByIDOut
	}{
		{
			name: "IsProjectAuthor",
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{
					Name:            "test",
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetMetaByIDOut{Name: "test"},
		},
		{
			name: "IsAuthor",
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{
					Name:            "test",
					AuthorID:        1,
					ProjectAuthorID: 2,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetMetaByIDOut{Name: "test"},
		},
		{
			name: "IsReader",
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{
					Name:            "test",
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleRead}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetMetaByIDOut{Name: "test"},
		},
		{
			name: "IsWriter",
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{
					Name:            "test",
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetMetaByIDOut{Name: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)

			tt.setup(templateRepo)

			usecase := New(templateRepo)

			got, err := usecase.Handle(ctx, in)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	in := domain.TemplateGetMetaByIDIn{TemplateID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{
					AuthorID:        2,
					ProjectAuthorID: 3,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)

			tt.setup(templateRepo)

			usecase := New(templateRepo)

			_, err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
