package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	var in domain.FolderCreateIn
	require.NoError(t, gofakeit.Struct(&in))

	tests := []struct {
		name  string
		setup func(folderRepo *MockfolderRepository)
	}{
		{
			name: "IsAuthor",
			setup: func(folderRepo *MockfolderRepository) {
				folder := domain.Folder{AuthorID: in.AuthorID}
				folderRepo.EXPECT().GetByID(ctx, in.ParentID).Return(&folder, nil)
				folderRepo.EXPECT().Create(ctx, in.ParentID, in.Name, in.AuthorID, folder.RootAuthorID).Return(nil)
			},
		},
		{
			name: "IsRootAuthor",
			setup: func(folderRepo *MockfolderRepository) {
				folder := domain.Folder{RootAuthorID: in.AuthorID}
				folderRepo.EXPECT().GetByID(ctx, in.ParentID).Return(&folder, nil)
				folderRepo.EXPECT().Create(ctx, in.ParentID, in.Name, in.AuthorID, folder.RootAuthorID).Return(nil)
			},
		},
		{
			name: "IsMaintainer",
			setup: func(folderRepo *MockfolderRepository) {
				folder := domain.Folder{
					Users: []domain.FolderUser{{ID: in.AuthorID, Role: user_domain.RoleMaintain}},
				}
				folderRepo.EXPECT().GetByID(ctx, in.ParentID).Return(&folder, nil)
				folderRepo.EXPECT().Create(ctx, in.ParentID, in.Name, in.AuthorID, folder.RootAuthorID).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			tt.setup(folderRepo)

			usecase := New(folderRepo)
			err := usecase.Handle(ctx, in)
			require.NoError(t, err)
		})
	}

}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	var in domain.FolderCreateIn
	require.NoError(t, gofakeit.Struct(&in))

	tests := []struct {
		name  string
		setup func(folderRepo *MockfolderRepository)
		in    domain.FolderCreateIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(folderRepo *MockfolderRepository) {},
			in:    domain.FolderCreateIn{Name: ""},
			want:  domain.ErrEmptyValue.Error(),
		},
		{
			name: "folderRepo_GetByID",
			setup: func(folderRepo *MockfolderRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   in,
			want: "test1",
		},
		{
			name: "domain_ErrParentNotFound",
			setup: func(folderRepo *MockfolderRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			in:   in,
			want: domain.ErrParentNotFound.Error(),
		},
		{
			name: "domain_ErrParentInvalid",
			setup: func(folderRepo *MockfolderRepository) {
				folder := domain.Folder{
					AuthorID:     gofakeit.Int64(),
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
			},
			in:   in,
			want: domain.ErrParentInvalid.Error(),
		},
		{
			name: "folderRepo_Create",
			setup: func(folderRepo *MockfolderRepository) {
				folder := domain.Folder{
					AuthorID:     in.AuthorID,
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folder, nil)
				folderRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test2"))
			},
			in:   in,
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			folderRepo := NewMockfolderRepository(ctrl)
			tt.setup(folderRepo)

			usecase := New(folderRepo)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
