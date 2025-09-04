package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockuserRepository(ctrl)
	hasher := NewMockpasswordHasher(ctrl)

	in := domain.UserCreateIn{
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: domain.Password("aB12345-"),
	}

	userRepo.EXPECT().ExistsByNameOrEmail(trCtx, in.Name, in.Email).Return(false, nil)
	hasher.EXPECT().Hash(in.Password).Return([]byte{1, 2, 3}, nil)
	userRepo.EXPECT().Create(trCtx, in.Name, in.Email, []byte{1, 2, 3}).Return(nil)

	usecase := New(userRepo, hasher)
	err := usecase.Handle(ctx, in)
	require.NoError(t, err)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	testErr := errors.New("test error")

	validIn := domain.UserCreateIn{
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: domain.Password("aB12345-"),
	}

	tests := []struct {
		name  string
		setup func(userRepo *MockuserRepository, hasher *MockpasswordHasher)
		in    domain.UserCreateIn
		want  error
	}{
		{
			name:  "in_Validate",
			setup: func(userRepo *MockuserRepository, hasher *MockpasswordHasher) {},
			in:    domain.UserCreateIn{},
			want:  domain.ErrEmptyValue,
		},
		{
			name: "userRepo_ExistsByNameOrEmail",
			setup: func(userRepo *MockuserRepository, hasher *MockpasswordHasher) {
				userRepo.EXPECT().ExistsByNameOrEmail(trCtx, gomock.Any(), gomock.Any()).Return(false, testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "domain_UserExists",
			setup: func(userRepo *MockuserRepository, hasher *MockpasswordHasher) {
				userRepo.EXPECT().ExistsByNameOrEmail(trCtx, gomock.Any(), gomock.Any()).Return(true, nil)
			},
			in:   validIn,
			want: domain.ErrUserExists,
		},
		{
			name: "hasher_Hash",
			setup: func(userRepo *MockuserRepository, hasher *MockpasswordHasher) {
				userRepo.EXPECT().ExistsByNameOrEmail(trCtx, gomock.Any(), gomock.Any()).Return(false, nil)
				hasher.EXPECT().Hash(gomock.Any()).Return(nil, testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "userRepo_Create",
			setup: func(userRepo *MockuserRepository, hasher *MockpasswordHasher) {
				userRepo.EXPECT().ExistsByNameOrEmail(trCtx, gomock.Any(), gomock.Any()).Return(false, nil)
				hasher.EXPECT().Hash(gomock.Any()).Return([]byte{}, nil)
				userRepo.EXPECT().Create(trCtx, gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
			},
			in:   validIn,
			want: testErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := NewMockuserRepository(ctrl)
			hasher := NewMockpasswordHasher(ctrl)
			tt.setup(userRepo, hasher)

			usecase := New(userRepo, hasher)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorIs(t, err, tt.want)
		})
	}
}
