package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	want := []domain.TemplateUser{
		{ID: 2, Name: "alice", Email: "alice@example.com", Role: user_domain.RoleRead},
		{ID: 3, Name: "bob", Email: "bob@example.com", Role: user_domain.RoleWrite},
	}

	tests := []struct {
		name  string
		in    domain.TemplateUserListIn
		setup func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository)
	}{
		{
			name: "TemplateAuthor",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&domain.Template{AuthorID: 1, ProjectAuthorID: 9}, nil)
				templateUserRepo.EXPECT().GetByTemplateID(ctx, int64(10)).Return(want, nil)
			},
		},
		{
			name: "ProjectAuthor",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&domain.Template{AuthorID: 9, ProjectAuthorID: 1}, nil)
				templateUserRepo.EXPECT().GetByTemplateID(ctx, int64(10)).Return(want, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateUserRepo := NewMocktemplateUserRepository(ctrl)
			tt.setup(templateRepo, templateUserRepo)

			usecase := New(templateRepo, templateUserRepo)
			got, err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.TemplateUserListIn
		setup func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Template{AuthorID: 8, ProjectAuthorID: 9}, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateUserRepo_GetByTemplateID",
			in:   domain.TemplateUserListIn{UserID: 1, TemplateID: 10},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Template{AuthorID: 1, ProjectAuthorID: 9}, nil)
				templateUserRepo.EXPECT().GetByTemplateID(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateUserRepo := NewMocktemplateUserRepository(ctrl)
			tt.setup(templateRepo, templateUserRepo)

			usecase := New(templateRepo, templateUserRepo)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
