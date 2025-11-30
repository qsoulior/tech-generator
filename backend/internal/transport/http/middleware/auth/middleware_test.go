package auth_middleware

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

func TestMiddleware_Handle_Success(t *testing.T) {
	ctx := context.Background()

	t.Run("Pass", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecase := NewMockusecase(ctrl)

		token := gofakeit.UUID()
		user := domain.User{ID: gofakeit.Int64()}
		usecase.EXPECT().Handle(ctx, token).Return(&user, nil)

		middleware := New(usecase, nil)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {
			userIDString := r.Header.Get(headerUserID)
			userID, err := strconv.ParseInt(userIDString, 10, 64)
			require.NoError(t, err)
			require.Equal(t, user.ID, userID)
		}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "token", Value: token})

		handlerWrapped.ServeHTTP(rr, r)
	})

	t.Run("Skip", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecase := NewMockusecase(ctrl)

		middleware := New(usecase, nil)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodPost, "/user/token/create", nil)

		handlerWrapped.ServeHTTP(rr, r)

		resp := rr.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestMiddleware_Handle_Error(t *testing.T) {
	ctx := context.Background()

	t.Run("NoCookie", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecase := NewMockusecase(ctrl)

		middleware := New(usecase, nil)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		handlerWrapped.ServeHTTP(rr, r)

		resp := rr.Result()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("CookieInvalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecase := NewMockusecase(ctrl)

		token := gofakeit.UUID()
		expectedErr := error_domain.NewBaseError("test")
		usecase.EXPECT().Handle(ctx, token).Return(nil, expectedErr)

		middleware := New(usecase, nil)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "token", Value: token})

		handlerWrapped.ServeHTTP(rr, r)

		resp := rr.Result()
		require.Equal(t, http.StatusForbidden, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, expectedErr.Error(), string(body))
	})

	t.Run("BaseError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecase := NewMockusecase(ctrl)

		token := gofakeit.UUID()
		expectedErr := error_domain.NewBaseError("test")
		usecase.EXPECT().Handle(ctx, token).Return(nil, expectedErr)

		middleware := New(usecase, nil)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "token", Value: token})

		handlerWrapped.ServeHTTP(rr, r)

		resp := rr.Result()
		require.Equal(t, http.StatusForbidden, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, expectedErr.Error(), string(body))
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		logHandler := slog.NewJSONHandler(&buf, nil)
		logger := slog.New(logHandler)

		usecase := NewMockusecase(ctrl)

		token := gofakeit.UUID()
		expectedErr := errors.New("test")
		usecase.EXPECT().Handle(ctx, token).Return(nil, expectedErr)

		middleware := New(usecase, logger)
		fn := middleware.Handle()

		handlerBase := func(w http.ResponseWriter, r *http.Request) {}

		handlerWrapped := fn(http.HandlerFunc(handlerBase))

		rr := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "token", Value: token})

		handlerWrapped.ServeHTTP(rr, r)

		resp := rr.Result()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
