package env

import (
	"os"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestGetString(t *testing.T) {
	t.Run("Exist", func(t *testing.T) {
		key := gofakeit.UUID()
		value := gofakeit.UUID()
		fallback := gofakeit.UUID()

		err := os.Setenv(key, value)
		require.NoError(t, err)
		defer func() { require.NoError(t, os.Unsetenv(key)) }()

		got := GetString(key, fallback)
		require.Equal(t, value, got)
	})

	t.Run("NotExist", func(t *testing.T) {
		key := gofakeit.UUID()
		fallback := gofakeit.UUID()
		got := GetString(key, fallback)
		require.Equal(t, fallback, got)
	})
}

func TestGetInteger(t *testing.T) {
	t.Run("Exist", func(t *testing.T) {
		key := gofakeit.UUID()
		value := gofakeit.Int()
		fallback := gofakeit.Int()

		err := os.Setenv(key, strconv.Itoa(value))
		require.NoError(t, err)
		defer func() { require.NoError(t, os.Unsetenv(key)) }()

		got := GetInteger(key, fallback)
		require.Equal(t, value, got)
	})

	t.Run("FailParse", func(t *testing.T) {
		key := gofakeit.UUID()
		value := gofakeit.UUID()
		fallback := gofakeit.Int()

		err := os.Setenv(key, value)
		require.NoError(t, err)
		defer func() { require.NoError(t, os.Unsetenv(key)) }()

		got := GetInteger(key, fallback)
		require.Equal(t, fallback, got)
	})

	t.Run("NotExist", func(t *testing.T) {
		key := gofakeit.UUID()
		fallback := gofakeit.Int()
		got := GetInteger(key, fallback)
		require.Equal(t, fallback, got)
	})
}

func TestGetFloat(t *testing.T) {
	t.Run("Exist", func(t *testing.T) {
		key := gofakeit.UUID()
		value := gofakeit.Float64()
		fallback := gofakeit.Float64()

		err := os.Setenv(key, strconv.FormatFloat(value, 'f', -1, 64))
		require.NoError(t, err)
		defer func() { require.NoError(t, os.Unsetenv(key)) }()

		got := GetFloat(key, fallback)
		require.Equal(t, value, got)
	})

	t.Run("FailParse", func(t *testing.T) {
		key := gofakeit.UUID()
		value := gofakeit.UUID()
		fallback := gofakeit.Float64()

		err := os.Setenv(key, value)
		require.NoError(t, err)
		defer func() { require.NoError(t, os.Unsetenv(key)) }()

		got := GetFloat(key, fallback)
		require.Equal(t, fallback, got)
	})

	t.Run("NotExist", func(t *testing.T) {
		key := gofakeit.UUID()
		fallback := gofakeit.Float64()
		got := GetFloat(key, fallback)
		require.Equal(t, fallback, got)
	})
}
