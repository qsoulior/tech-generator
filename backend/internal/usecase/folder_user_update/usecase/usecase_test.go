package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.FolderUserUpdateIn
		setup func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users: []domain.FolderUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
					{ID: 4, Role: user_domain.RoleRead},  // old user
					{ID: 6, Role: user_domain.RoleRead},  // not existing user
				},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 1, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, int64(10)).Return(&folder, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, []int64{2, 3, 4, 6}).Return(userExistingIDs, nil)

				folderUsersExisting := []domain.FolderUser{
					{ID: 3, Role: user_domain.RoleRead}, // new user role
					{ID: 4, Role: user_domain.RoleRead}, // old user
					{ID: 5, Role: user_domain.RoleRead}, // missing user
				}
				folderUserRepo.EXPECT().GetByFolderID(ctx, int64(10)).Return(folderUsersExisting, nil)

				folderUsersToUpsert := []domain.FolderUser{
					{ID: 2, Role: user_domain.RoleRead},  // new user
					{ID: 3, Role: user_domain.RoleWrite}, // new user role
				}
				folderUserRepo.EXPECT().Upsert(trCtx, int64(10), folderUsersToUpsert).Return(nil)

				folderUserIDsToDelete := []int64{5}
				folderUserRepo.EXPECT().Delete(trCtx, int64(10), folderUserIDsToDelete).Return(nil)
			},
		},
		{
			name: "IsRootAuthor/NoUpdates",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 9, RootAuthorID: 1}
				folderRepo.EXPECT().GetByID(ctx, int64(10)).Return(&folder, nil)

				userRepo.EXPECT().GetByIDs(ctx, []int64{2}).Return([]int64{2}, nil)

				folderUsersExisting := []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}}
				folderUserRepo.EXPECT().GetByFolderID(ctx, int64(10)).Return(folderUsersExisting, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			folderUserRepo := NewMockfolderUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(folderRepo, folderUserRepo, userRepo)

			usecase := New(folderRepo, folderUserRepo, userRepo, trManager)
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
		in    domain.FolderUserUpdateIn
		setup func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository)
		want  string
	}{
		{
			name: "folderRepo_GetByID",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrFolderNotFound",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrFolderNotFound.Error(),
		},
		{
			name: "domain_ErrFolderInvalid",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 8, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
			},
		},
		{
			name: "userRepo_GetByIDs",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 1, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "folderUserRepo_GetByFolderID",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 1, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)
				folderUserRepo.EXPECT().GetByFolderID(ctx, gomock.Any()).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "folderUserRepo_Upsert",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 1, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return([]int64{2}, nil)

				folderUsersExisting := []domain.FolderUser{{ID: 3, Role: user_domain.RoleRead}}
				folderUserRepo.EXPECT().GetByFolderID(ctx, gomock.Any()).Return(folderUsersExisting, nil)
				folderUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "folderUserRepo_Delete",
			in: domain.FolderUserUpdateIn{
				UserID:   1,
				FolderID: 10,
				Users:    []domain.FolderUser{{ID: 2, Role: user_domain.RoleRead}},
			},
			setup: func(folderRepo *MockfolderRepository, folderUserRepo *MockfolderUserRepository, userRepo *MockuserRepository) {
				folder := domain.Folder{AuthorID: 1, RootAuthorID: 9}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)

				userExistingIDs := []int64{2, 3, 4}
				userRepo.EXPECT().GetByIDs(ctx, gomock.Any()).Return(userExistingIDs, nil)

				folderUsersExisting := []domain.FolderUser{{ID: 3, Role: user_domain.RoleRead}}
				folderUserRepo.EXPECT().GetByFolderID(ctx, gomock.Any()).Return(folderUsersExisting, nil)
				folderUserRepo.EXPECT().Upsert(trCtx, gomock.Any(), gomock.Any()).Return(nil)
				folderUserRepo.EXPECT().Delete(trCtx, gomock.Any(), gomock.Any()).Return(errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			folderUserRepo := NewMockfolderUserRepository(ctrl)
			userRepo := NewMockuserRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(folderRepo, folderUserRepo, userRepo)

			usecase := New(folderRepo, folderUserRepo, userRepo, trManager)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
