package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.ProjectUserUpdateIn
		setup func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users: []domain.ProjectUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
					{ID: 4, Role: user_domain.RoleRead},  // old user
					{ID: 6, Role: user_domain.RoleRead},  // not existing user
				},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, []int64{2, 3, 4, 6}).Return(userExistingIDs, nil)

				projectUsersExisting := []domain.ProjectUser{
					{ID: 3, Role: user_domain.RoleRead}, // new user role
					{ID: 4, Role: user_domain.RoleRead}, // old user
					{ID: 5, Role: user_domain.RoleRead}, // missing user
				}
				projectUserRepo.EXPECT().GetByProjectID(ctx, int64(10)).Return(projectUsersExisting, nil)

				projectUsersToUpsert := []domain.ProjectUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
				}
				projectUserRepo.EXPECT().Upsert(trCtx, int64(10), projectUsersToUpsert).Return(nil)

				projectUserIDsToDelete := []int64{5}
				projectUserRepo.EXPECT().Delete(trCtx, int64(10), projectUserIDsToDelete).Return(nil)
			},
		},
		{
			name: "NoUpdates",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)

				userRepo.EXPECT().GetByIDs(ctx, []int64{2}).Return([]int64{2}, nil)

				projectUsersExisting := []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}}
				projectUserRepo.EXPECT().GetByProjectID(ctx, int64(10)).Return(projectUsersExisting, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			projectUserRepo := NewMockprojectUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(projectRepo, projectUserRepo, userRepo)

			usecase := New(projectRepo, projectUserRepo, userRepo, trManager)
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
		in    domain.ProjectUserUpdateIn
		setup func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository)
		want  string
	}{
		{
			name: "projectRepo_GetByID",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 8}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
			},
		},
		{
			name: "userRepo_GetByIDs",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "projectUserRepo_GetByProjectID",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)
				projectUserRepo.EXPECT().GetByProjectID(ctx, gomock.Any()).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "projectUserRepo_Upsert",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)

				projectUsersExisting := []domain.ProjectUser{{ID: 3, Role: user_domain.RoleRead}}
				projectUserRepo.EXPECT().GetByProjectID(ctx, gomock.Any()).Return(projectUsersExisting, nil)
				projectUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "projectUserRepo_Delete",
			in: domain.ProjectUserUpdateIn{
				UserID:    1,
				ProjectID: 10,
				Users:     []domain.ProjectUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository, userRepo *MockuserRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&project, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(userExistingIDs, nil)

				projectUsersExisting := []domain.ProjectUser{{ID: 3, Role: user_domain.RoleRead}}
				projectUserRepo.EXPECT().GetByProjectID(ctx, gomock.Any()).Return(projectUsersExisting, nil)
				projectUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(nil)
				projectUserRepo.EXPECT().Delete(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			projectUserRepo := NewMockprojectUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(projectRepo, projectUserRepo, userRepo)

			usecase := New(projectRepo, projectUserRepo, userRepo, trManager)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
