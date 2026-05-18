package user_token_delete_handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_UserTokenDelete_ClearsCookie(t *testing.T) {
	handler := New()
	got, err := handler.UserTokenDelete(context.Background())
	require.NoError(t, err)
	require.NotNil(t, got)

	cookie, err := http.ParseSetCookie(got.SetCookie)
	require.NoError(t, err)
	require.Equal(t, "token", cookie.Name)
	require.Empty(t, cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.HttpOnly)
	require.True(t, cookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Negative(t, cookie.MaxAge, "delete cookie must carry negative Max-Age; got %d", cookie.MaxAge)
}
