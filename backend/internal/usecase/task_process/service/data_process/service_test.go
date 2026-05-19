package data_process_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
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

func TestService_Handle_Sprig(t *testing.T) {
	ctx := context.Background()
	service := New()

	tests := []struct {
		name string
		in   domain.DataProcessIn
		want string
	}{
		{
			name: "printf_decimals",
			in: domain.DataProcessIn{
				Values: map[string]any{"price": 1234.5678},
				Data:   []byte(`{{ printf "%.2f" .price }}`),
			},
			want: "1234.57",
		},
		{
			name: "default_fallback",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte(`{{ default "n/a" .missing }}`),
			},
			want: "n/a",
		},
		{
			name: "upper_replace",
			in: domain.DataProcessIn{
				Values: map[string]any{"name": "foo bar"},
				Data:   []byte(`{{ .name | replace " " "_" | upper }}`),
			},
			want: "FOO_BAR",
		},
		{
			name: "range_list",
			in: domain.DataProcessIn{
				Values: map[string]any{"items": []string{"a", "b", "c"}},
				Data:   []byte(`{{- range .items }}{{ . }}{{ end -}}`),
			},
			want: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Handle(ctx, tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
		})
	}
}

func TestService_Handle_Error(t *testing.T) {
	ctx := context.Background()
	service := New()

	tests := []struct {
		name        string
		in          domain.DataProcessIn
		wantMessage string
		wantLine    int
		wantSnippet string
	}{
		{
			name: "TemplateParse",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte("abc {{abc}}"),
			},
			wantMessage: task_domain.MessageTemplateParse,
			wantLine:    1,
			wantSnippet: "abc {{abc}}",
		},
		{
			name: "EnvBlocked",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte(`{{ env "PATH" }}`),
			},
			wantMessage: task_domain.MessageTemplateParse,
			wantLine:    1,
			wantSnippet: `{{ env "PATH" }}`,
		},
		{
			name: "ExpandEnvBlocked",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte(`{{ expandenv "$PATH" }}`),
			},
			wantMessage: task_domain.MessageTemplateParse,
			wantLine:    1,
			wantSnippet: `{{ expandenv "$PATH" }}`,
		},
		{
			name: "TemplateParseOnSecondLine",
			in: domain.DataProcessIn{
				Values: map[string]any{},
				Data:   []byte("first line\nbroken {{abc}}\nthird line"),
			},
			wantMessage: task_domain.MessageTemplateParse,
			wantLine:    2,
			wantSnippet: "broken {{abc}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Handle(ctx, tt.in)

			var got *task_domain.ProcessError
			require.ErrorAs(t, err, &got)
			require.Equal(t, tt.wantMessage, got.Message)
			require.NotNil(t, got.Template, "expected template error to be populated")
			require.Equal(t, tt.wantLine, got.Template.Line)
			require.Equal(t, tt.wantSnippet, got.Template.Snippet)
			require.NotEmpty(t, got.Template.Detail)
		})
	}
}
