package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockuserRepository(ctrl)
	passwordVerifier := NewMockpasswordVerifier(ctrl)
	tokenBuilder := NewMocktokenBuilder(ctrl)
	tokenExpiration := 7 * 24 * time.Hour

	in := domain.UserCreateTokenIn{
		Name:     gofakeit.Username(),
		Password: domain.Password("aB12345-"),
	}

	user := domain.User{
		ID:       gofakeit.Int64(),
		Password: []byte{1, 2, 3},
	}
	userRepo.EXPECT().GetByName(ctx, in.Name).Return(&user, nil)

	passwordVerifier.EXPECT().Verify(user.Password, in.Password).Return(nil)

	var want domain.UserCreateTokenOut
	require.NoError(t, gofakeit.Struct(&want))
	tokenBuilder.EXPECT().Build(user.ID, tokenExpiration).Return(want, nil)

	usecase := New(userRepo, passwordVerifier, tokenBuilder, tokenExpiration)
	got, err := usecase.Handle(ctx, in)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	testErr := errors.New("test error")

	validIn := domain.UserCreateTokenIn{
		Name:     gofakeit.Username(),
		Password: domain.Password("aB12345-"),
	}

	tests := []struct {
		name  string
		setup func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder)
		in    domain.UserCreateTokenIn
		want  error
	}{
		{
			name: "in_Validate",
			setup: func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder) {
			},
			in:   domain.UserCreateTokenIn{},
			want: domain.ErrNameEmpty,
		},
		{
			name: "userRepo_GetByName",
			setup: func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder) {
				userRepo.EXPECT().GetByName(ctx, gomock.Any()).Return(nil, testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "userRepo_GetByName_UserDoesNotExist",
			setup: func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder) {
				userRepo.EXPECT().GetByName(ctx, gomock.Any()).Return(nil, nil)
			},
			in:   validIn,
			want: domain.ErrUserDoesNotExist,
		},
		{
			name: "passwordVerifier_Verify_PasswordIncorrect",
			setup: func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder) {
				user := domain.User{}
				userRepo.EXPECT().GetByName(ctx, gomock.Any()).Return(&user, nil)
				passwordVerifier.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(testErr)
			},
			in:   validIn,
			want: domain.ErrPasswordIncorrect,
		},
		{
			name: "tokenBuilder_Build",
			setup: func(userRepo *MockuserRepository, passwordVerifier *MockpasswordVerifier, tokenBuilder *MocktokenBuilder) {
				user := domain.User{}
				userRepo.EXPECT().GetByName(ctx, gomock.Any()).Return(&user, nil)
				passwordVerifier.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(nil)
				tokenBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(domain.UserCreateTokenOut{}, testErr)
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
			passwordVerifier := NewMockpasswordVerifier(ctrl)
			tokenBuilder := NewMocktokenBuilder(ctrl)
			tt.setup(userRepo, passwordVerifier, tokenBuilder)

			usecase := New(userRepo, passwordVerifier, tokenBuilder, 0)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorIs(t, err, tt.want)
		})
	}
}
