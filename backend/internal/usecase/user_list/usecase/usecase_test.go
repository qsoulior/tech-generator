package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	in := domain.UserListIn{
		Page:   1,
		Size:   10,
		Filter: domain.UserListFilter{ExcludeUserID: 1},
	}
	users := []domain.User{
		{ID: 2, Name: "alice", Email: "alice@example.com"},
		{ID: 3, Name: "bob", Email: "bob@example.com"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockuserRepository(ctrl)
	userRepo.EXPECT().List(ctx, in).Return(users, nil)
	userRepo.EXPECT().GetTotal(ctx, in).Return(int64(11), nil)

	usecase := New(userRepo)
	got, err := usecase.Handle(ctx, in)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, users, got.Users)
	require.Equal(t, int64(11), got.TotalUsers)
	require.Equal(t, int64(2), got.TotalPages)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.UserListIn
		setup func(userRepo *MockuserRepository)
		want  string
	}{
		{
			name: "PageInvalid",
			in:   domain.UserListIn{Page: 0, Size: 10},
			setup: func(userRepo *MockuserRepository) {
			},
			want: "page",
		},
		{
			name: "SizeInvalid",
			in:   domain.UserListIn{Page: 1, Size: 0},
			setup: func(userRepo *MockuserRepository) {
			},
			want: "size",
		},
		{
			name: "userRepo_List",
			in:   domain.UserListIn{Page: 1, Size: 10},
			setup: func(userRepo *MockuserRepository) {
				userRepo.EXPECT().List(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "userRepo_GetTotal",
			in:   domain.UserListIn{Page: 1, Size: 10},
			setup: func(userRepo *MockuserRepository) {
				userRepo.EXPECT().List(ctx, gomock.Any()).Return([]domain.User{}, nil)
				userRepo.EXPECT().GetTotal(ctx, gomock.Any()).Return(int64(0), errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := NewMockuserRepository(ctrl)
			tt.setup(userRepo)

			usecase := New(userRepo)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
