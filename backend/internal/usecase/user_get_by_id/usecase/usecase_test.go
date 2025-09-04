package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var want domain.User
	require.NoError(t, gofakeit.Struct(&want))

	userRepo := NewMockuserRepository(ctrl)
	userRepo.EXPECT().GetByID(ctx, want.ID).Return(&want, nil)

	usecase := New(userRepo)
	got, err := usecase.Handle(ctx, want.ID)
	require.NoError(t, err)
	require.Equal(t, want, *got)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := errors.New("want err")

	userRepo := NewMockuserRepository(ctrl)
	userRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, want)

	usecase := New(userRepo)
	_, err := usecase.Handle(ctx, gofakeit.Int64())
	require.ErrorIs(t, err, want)
}
