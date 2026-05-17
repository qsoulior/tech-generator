package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockuserRepository(ctrl)
	passwordHasher := NewMockpasswordHasher(ctrl)

	in := domain.UserCreateIn{
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: domain.Password("aB12345-"),
	}

	passwordHasher.EXPECT().Hash(in.Password).Return([]byte{1, 2, 3}, nil)
	userRepo.EXPECT().Create(ctx, in.Name, in.Email, []byte{1, 2, 3}).Return(nil)

	usecase := New(userRepo, passwordHasher)
	err := usecase.Handle(ctx, in)
	require.NoError(t, err)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	testErr := errors.New("test error")

	validIn := domain.UserCreateIn{
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: domain.Password("aB12345-"),
	}

	tests := []struct {
		name  string
		setup func(userRepo *MockuserRepository, passwordHasher *MockpasswordHasher)
		in    domain.UserCreateIn
		want  error
	}{
		{
			name:  "in_Validate",
			setup: func(userRepo *MockuserRepository, passwordHasher *MockpasswordHasher) {},
			in:    domain.UserCreateIn{},
			want:  domain.ErrValueEmpty,
		},
		{
			name: "passwordHasher_Hash",
			setup: func(userRepo *MockuserRepository, passwordHasher *MockpasswordHasher) {
				passwordHasher.EXPECT().Hash(gomock.Any()).Return(nil, testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "userRepo_Create",
			setup: func(userRepo *MockuserRepository, passwordHasher *MockpasswordHasher) {
				passwordHasher.EXPECT().Hash(gomock.Any()).Return([]byte{}, nil)
				userRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "userRepo_Create_UserExists",
			setup: func(userRepo *MockuserRepository, passwordHasher *MockpasswordHasher) {
				passwordHasher.EXPECT().Hash(gomock.Any()).Return([]byte{}, nil)
				userRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrUserExists)
			},
			in:   validIn,
			want: domain.ErrUserExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := NewMockuserRepository(ctrl)
			passwordHasher := NewMockpasswordHasher(ctrl)
			tt.setup(userRepo, passwordHasher)

			usecase := New(userRepo, passwordHasher)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorIs(t, err, tt.want)
		})
	}
}
