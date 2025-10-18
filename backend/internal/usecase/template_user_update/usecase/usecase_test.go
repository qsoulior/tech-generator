package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.TemplateUserUpdateIn
		setup func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users: []domain.TemplateUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
					{ID: 4, Role: user_domain.RoleRead},  // old user
					{ID: 6, Role: user_domain.RoleRead},  // not existing user
				},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 1, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, []int64{2, 3, 4, 6}).Return(userExistingIDs, nil)

				templateUsersExisting := []domain.TemplateUser{
					{ID: 3, Role: user_domain.RoleRead}, // new user role
					{ID: 4, Role: user_domain.RoleRead}, // old user
					{ID: 5, Role: user_domain.RoleRead}, // missing user
				}
				templateUserRepo.EXPECT().GetByTemplateID(ctx, int64(10)).Return(templateUsersExisting, nil)

				templateUsersToUpsert := []domain.TemplateUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
				}
				templateUserRepo.EXPECT().Upsert(trCtx, int64(10), templateUsersToUpsert).Return(nil)

				templateUserIDsToDelete := []int64{5}
				templateUserRepo.EXPECT().Delete(trCtx, int64(10), templateUserIDsToDelete).Return(nil)
			},
		},
		{
			name: "IsRootAuthor/NoUpdates",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 9, RootAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				userRepo.EXPECT().GetByIDs(ctx, []int64{2}).Return([]int64{2}, nil)

				templateUsersExisting := []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}}
				templateUserRepo.EXPECT().GetByTemplateID(ctx, int64(10)).Return(templateUsersExisting, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateUserRepo := NewMocktemplateUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, templateUserRepo, userRepo)

			usecase := New(templateRepo, templateUserRepo, userRepo, trManager)
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.TemplateUserUpdateIn
		setup func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 8, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&template, nil)
			},
		},
		{
			name: "userRepo_GetByIDs",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 1, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&template, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "templateUserRepo_GetByTemplateID",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 1, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&template, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)
				templateUserRepo.EXPECT().GetByTemplateID(ctx, gomock.Any()).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "templateUserRepo_Upsert",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 1, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&template, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)

				templateUsersExisting := []domain.TemplateUser{{ID: 3, Role: user_domain.RoleRead}}
				templateUserRepo.EXPECT().GetByTemplateID(ctx, gomock.Any()).Return(templateUsersExisting, nil)
				templateUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "templateUserRepo_Delete",
			in: domain.TemplateUserUpdateIn{
				UserID:     1,
				TemplateID: 10,
				Users:      []domain.TemplateUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(templateRepo *MocktemplateRepository, templateUserRepo *MocktemplateUserRepository, userRepo *MockuserRepository) {
				template := domain.Template{AuthorID: 1, RootAuthorID: 9}
				templateRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&template, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(userExistingIDs, nil)

				templateUsersExisting := []domain.TemplateUser{{ID: 3, Role: user_domain.RoleRead}}
				templateUserRepo.EXPECT().GetByTemplateID(ctx, gomock.Any()).Return(templateUsersExisting, nil)
				templateUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(nil)
				templateUserRepo.EXPECT().Delete(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateUserRepo := NewMocktemplateUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, templateUserRepo, userRepo)

			usecase := New(templateRepo, templateUserRepo, userRepo, trManager)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
