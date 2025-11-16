package data_process_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

func TestService_Handle_Success(t *testing.T) {
	ctx := context.Background()
	service := New()

	in := domain.DataProcessIn{
		Values: map[string]any{
			"test1": 123,
			"test2": "some text",
			"test3": 456.789,
			"test4": "foo bar",
		},
		Data: []byte("Lorem {{.test1}} ipsum {{ .test3 }}, {{.test4}}"),
	}

	got, err := service.Handle(ctx, in)
	require.NoError(t, err)

	want := "Lorem 123 ipsum 456.789, foo bar"
	require.Equal(t, want, string(got))
}

func TestService_Handle_Error(t *testing.T) {
	ctx := context.Background()
	service := New()

	tests := []struct {
		name string
		in   domain.DataProcessIn
		want domain.ProcessError
	}{
		{
			name: "TemplateParse",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte("abc {{abc}}"),
			},
			want: domain.ProcessError{Message: domain.MessageTemplateParse},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Handle(ctx, tt.in)

			var got *domain.ProcessError
			require.ErrorAs(t, err, &got)
			require.Equal(t, tt.want, *got)
		})
	}
}
