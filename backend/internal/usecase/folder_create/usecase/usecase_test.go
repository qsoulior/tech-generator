package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.FolderCreateIn
		setup func(folderRepo *MockfolderRepository)
	}{
		{
			name: "WithoutParent",
			in: domain.FolderCreateIn{
				ParentID: nil,
				Name:     "test",
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository) {
				folderToCreate := domain.FolderToCreate{
					ParentID:     nil,
					Name:         "test",
					AuthorID:     1,
					RootAuthorID: 1,
				}
				folderRepo.EXPECT().Create(ctx, folderToCreate).Return(nil)
			},
		},
		{
			name: "WithParent_IsAuthor",
			in: domain.FolderCreateIn{
				ParentID: lo.ToPtr[int64](2),
				Name:     "test",
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository) {
				folderParent := domain.Folder{
					AuthorID:     1,
					RootAuthorID: 3,
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folderParent, nil)

				folderToCreate := domain.FolderToCreate{
					ParentID:     lo.ToPtr[int64](2),
					Name:         "test",
					AuthorID:     1,
					RootAuthorID: 3,
				}
				folderRepo.EXPECT().Create(ctx, folderToCreate).Return(nil)
			},
		},
		{
			name: "WithParent_IsRootAuthor",
			in: domain.FolderCreateIn{
				ParentID: lo.ToPtr[int64](2),
				Name:     "test",
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository) {
				folderParent := domain.Folder{
					AuthorID:     3,
					RootAuthorID: 1,
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folderParent, nil)

				folderToCreate := domain.FolderToCreate{
					ParentID:     lo.ToPtr[int64](2),
					Name:         "test",
					AuthorID:     1,
					RootAuthorID: 1,
				}
				folderRepo.EXPECT().Create(ctx, folderToCreate).Return(nil)
			},
		},
		{
			name: "WithParent_IsMaintainer",
			in: domain.FolderCreateIn{
				ParentID: lo.ToPtr[int64](2),
				Name:     "test",
				AuthorID: 1,
			},
			setup: func(folderRepo *MockfolderRepository) {
				folderParent := domain.Folder{
					AuthorID:     3,
					RootAuthorID: 4,
					Users:        []domain.FolderUser{{ID: 1, Role: user_domain.RoleMaintain}},
				}
				folderRepo.EXPECT().GetByID(ctx, int64(2)).Return(&folderParent, nil)

				folderToCreate := domain.FolderToCreate{
					ParentID:     lo.ToPtr[int64](2),
					Name:         "test",
					AuthorID:     1,
					RootAuthorID: 4,
				}
				folderRepo.EXPECT().Create(ctx, folderToCreate).Return(nil)
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
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

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
			in:   domain.FolderCreateIn{Name: "test", AuthorID: 1, ParentID: lo.ToPtr[int64](2)},
			want: "test1",
		},
		{
			name: "domain_ErrParentNotFound",
			setup: func(folderRepo *MockfolderRepository) {
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			in:   domain.FolderCreateIn{Name: "test", AuthorID: 1, ParentID: lo.ToPtr[int64](2)},
			want: domain.ErrParentNotFound.Error(),
		},
		{
			name: "domain_ErrParentInvalid",
			setup: func(folderRepo *MockfolderRepository) {
				folderParent := domain.Folder{
					AuthorID:     gofakeit.Int64(),
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folderParent, nil)
			},
			in:   domain.FolderCreateIn{Name: "test", AuthorID: 1, ParentID: lo.ToPtr[int64](2)},
			want: domain.ErrParentInvalid.Error(),
		},
		{
			name: "folderRepo_Create#1",
			setup: func(folderRepo *MockfolderRepository) {
				folderParent := domain.Folder{
					AuthorID:     1,
					RootAuthorID: gofakeit.Int64(),
					Users:        []domain.FolderUser{},
				}
				folderRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&folderParent, nil)
				folderRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test2"))
			},
			in:   domain.FolderCreateIn{Name: "test", AuthorID: 1, ParentID: lo.ToPtr[int64](2)},
			want: "test2",
		},
		{
			name: "folderRepo_Create#2",
			setup: func(folderRepo *MockfolderRepository) {
				folderRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test3"))
			},
			in:   domain.FolderCreateIn{Name: "test", AuthorID: 1, ParentID: nil},
			want: "test3",
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
