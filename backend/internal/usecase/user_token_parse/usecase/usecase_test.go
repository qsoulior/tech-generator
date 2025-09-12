package usecase

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tokenString := gofakeit.UUID()
	wantUser := domain.User{ID: gofakeit.Int64()}

	tokenParser := NewMocktokenParser(ctrl)
	tokenParser.EXPECT().Parse(tokenString).Return(&wantUser, nil)

	usecase := New(tokenParser)
	gotUser, err := usecase.Handle(ctx, tokenString)
	require.NoError(t, err)
	require.NotNil(t, gotUser)
	require.Equal(t, wantUser, *gotUser)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tokenString := gofakeit.UUID()
	expectedErr := gofakeit.Error()

	tokenParser := NewMocktokenParser(ctrl)
	tokenParser.EXPECT().Parse(tokenString).Return(nil, expectedErr)

	usecase := New(tokenParser)
	_, err := usecase.Handle(ctx, tokenString)
	require.ErrorIs(t, err, domain.ErrTokenInvalid)
}
