package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_update/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateUpdateIn
		setup func(templateRepo *MocktemplateRepository)
	}{
		{
			name: "IsTemplateAuthor",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(nil)
			},
		},
		{
			name: "IsProjectAuthor",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 2, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(nil)
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
		name    string
		in      domain.TemplateUpdateIn
		setup   func(templateRepo *MocktemplateRepository)
		want    string
		wantVal bool
	}{
		{
			name:    "ValidationEmptyName",
			in:      domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: ""},
			setup:   func(_ *MocktemplateRepository) {},
			want:    "name",
			wantVal: true,
		},
		{
			name: "templateRepo_GetByID",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 3, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateRepo_UpdateByID",
			in:   domain.TemplateUpdateIn{TemplateID: 10, UserID: 1, Name: "new"},
			setup: func(templateRepo *MocktemplateRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(errors.New("test2"))
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

			if tt.wantVal {
				var validationErr *error_domain.ValidationError
				require.ErrorAs(t, err, &validationErr)
			}
		})
	}
}
