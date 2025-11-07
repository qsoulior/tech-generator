package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateDeleteIn
		setup func(templateRepo *MocktemplateRepository)
	}{
		{
			name: "IsAuthor",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().DeleteByID(ctx, int64(10)).Return(nil)
			},
		},
		{
			name: "IsRootAuthor",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().DeleteByID(ctx, int64(10)).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(templateRepo)

			usecase := New(templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateDeleteIn
		setup func(templateRepo *MocktemplateRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateRepo_DeleteByID",
			in:   domain.TemplateDeleteIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().DeleteByID(ctx, int64(10)).Return(errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(templateRepo)

			usecase := New(templateRepo)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
