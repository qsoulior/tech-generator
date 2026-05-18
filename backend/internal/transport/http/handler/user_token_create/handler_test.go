package user_token_create_handler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

func TestHandler_UserTokenCreate_Success_Remember(t *testing.T) {
	ctx := context.Background()
	req := &api.UserTokenCreateRequest{Name: "alice", Password: "secret", Remember: true}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expiresAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	out := domain.UserCreateTokenOut{Token: "jwt", ExpiresAt: expiresAt}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.UserCreateTokenIn{Name: "alice", Password: domain.Password("secret")}).
		Return(out, nil)

	handler := New(usecase)
	got, err := handler.UserTokenCreate(ctx, req)
	require.NoError(t, err)

	resp, ok := got.(*api.UserTokenCreateCreated)
	require.True(t, ok, "expected *api.UserTokenCreateCreated, got %T", got)

	cookie, err := http.ParseSetCookie(resp.SetCookie)
	require.NoError(t, err)
	require.Equal(t, "token", cookie.Name)
	require.Equal(t, "jwt", cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.HttpOnly)
	require.True(t, cookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.WithinDuration(t, expiresAt, cookie.Expires, time.Second)
}

func TestHandler_UserTokenCreate_Success_SessionCookie(t *testing.T) {
	ctx := context.Background()
	req := &api.UserTokenCreateRequest{Name: "alice", Password: "secret", Remember: false}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	out := domain.UserCreateTokenOut{Token: "jwt", ExpiresAt: time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(out, nil)

	handler := New(usecase)
	got, err := handler.UserTokenCreate(ctx, req)
	require.NoError(t, err)

	resp, ok := got.(*api.UserTokenCreateCreated)
	require.True(t, ok, "expected *api.UserTokenCreateCreated, got %T", got)

	cookie, err := http.ParseSetCookie(resp.SetCookie)
	require.NoError(t, err)
	require.Equal(t, "token", cookie.Name)
	require.Equal(t, "jwt", cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.HttpOnly)
	require.True(t, cookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.True(t, cookie.Expires.IsZero(), "session cookie must not carry Expires; got %v", cookie.Expires)
	require.Zero(t, cookie.MaxAge, "session cookie must not carry Max-Age; got %d", cookie.MaxAge)
}

func TestHandler_UserTokenCreate_BaseError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserTokenCreateRequest{Name: "alice", Password: "secret", Remember: true}

	tests := []struct {
		name string
		err  error
	}{
		{name: "PasswordIncorrect", err: domain.ErrPasswordIncorrect},
		{name: "NameEmpty", err: domain.ErrNameEmpty},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(domain.UserCreateTokenOut{}, tt.err)

			handler := New(usecase)
			got, err := handler.UserTokenCreate(ctx, req)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_UserTokenCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserTokenCreateRequest{Name: "alice", Password: "", Remember: true}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("password", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(domain.UserCreateTokenOut{}, validationErr)

	handler := New(usecase)
	got, err := handler.UserTokenCreate(ctx, req)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_UserTokenCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.UserTokenCreateRequest{Name: "alice", Password: "secret", Remember: true}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(domain.UserCreateTokenOut{}, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.UserTokenCreate(ctx, req)
	require.Nil(t, got)
	require.ErrorContains(t, err, "user token create usecase")
	require.ErrorContains(t, err, "boom")
}
