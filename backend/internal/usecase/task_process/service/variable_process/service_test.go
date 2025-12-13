package variable_process_service

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

func TestService_Handle_Success(t *testing.T) {
	ctx := context.Background()
	service := New()

	in := domain.VariableProcessIn{
		Variables: []domain.Variable{
			// input integer variable with enabled constraints
			{
				ID:      1,
				Name:    "test1",
				Type:    variable_domain.TypeInteger,
				IsInput: true,
				Constraints: []domain.Constraint{
					{
						ID:         1,
						Name:       "test1",
						Expression: "test1 > 100",
						IsActive:   true,
					},
					{
						ID:         2,
						Name:       "test2",
						Expression: "test1 < 200",
						IsActive:   true,
					},
				},
			},
			// computed integer variable with disable constraint
			{
				ID:         2,
				Name:       "test2",
				Type:       variable_domain.TypeInteger,
				Expression: lo.ToPtr("test1 + test3 - test1 / 2"),
				IsInput:    false,
				Constraints: []domain.Constraint{
					{
						ID:         3,
						Name:       "test1",
						Expression: "test2 > 1000",
						IsActive:   false,
					},
				},
			},
			// input integer variable without constraints
			{
				ID:      3,
				Name:    "test3",
				Type:    variable_domain.TypeInteger,
				IsInput: true,
			},
			// input string variable without constraints
			{
				ID:      4,
				Name:    "test4",
				Type:    variable_domain.TypeString,
				IsInput: true,
			},
			// computed string variable with enabled constraints
			{
				ID:         5,
				Name:       "test5",
				Type:       variable_domain.TypeString,
				Expression: lo.ToPtr(`test4 + "bar"`),
				IsInput:    false,
				Constraints: []domain.Constraint{
					{
						ID:         4,
						Name:       "test1",
						Expression: "len(test5) > 5",
						IsActive:   true,
					},
				},
			},
			// computed float variable without constraints
			{
				ID:         6,
				Name:       "test6",
				Type:       variable_domain.TypeFloat,
				Expression: lo.ToPtr(`test2 + 10.2`),
				IsInput:    false,
			},
		},
		Payload: map[string]string{
			"test1": "123",
			"test3": "321",
			"test4": "foo",
		},
	}

	got, err := service.Handle(ctx, in)
	require.NoError(t, err)

	want := map[string]any{
		"test1": int64(123),
		"test2": 382.5,
		"test3": int64(321),
		"test4": "foo",
		"test5": "foobar",
		"test6": 392.7,
	}

	require.Equal(t, want, got)
}

func TestService_Handle_Error(t *testing.T) {
	ctx := context.Background()
	service := New()

	tests := []struct {
		name string
		in   domain.VariableProcessIn
		want task_domain.ProcessError
	}{
		{
			name: "ProcessError_Cycle",
			in: domain.VariableProcessIn{
				Variables: []domain.Variable{
					{
						ID:         1,
						Name:       "var1",
						Type:       variable_domain.TypeInteger,
						Expression: lo.ToPtr("var2 + 10"),
						IsInput:    false,
					},
					{
						ID:         2,
						Name:       "var2",
						Type:       variable_domain.TypeInteger,
						Expression: lo.ToPtr("var1 + 10"),
						IsInput:    false,
					},
				},
			},
			want: task_domain.ProcessError{
				Message: task_domain.MessageCycle,
			},
		},
		{
			name: "VariableError_Compile",
			in: domain.VariableProcessIn{
				Variables: []domain.Variable{
					{
						ID:         1,
						Name:       "var1",
						Type:       variable_domain.TypeInteger,
						Expression: lo.ToPtr("invalid1"),
						IsInput:    false,
					},
					{
						ID:         2,
						Name:       "var2",
						Type:       variable_domain.TypeInteger,
						Expression: lo.ToPtr("var1 + invalid"),
						IsInput:    false,
					},
				},
			},
			want: task_domain.ProcessError{
				VariableErrors: []task_domain.VariableError{
					{
						ID:      1,
						Name:    "var1",
						Message: task_domain.MessageVariableCompile,
					},
					{
						ID:      2,
						Name:    "var2",
						Message: task_domain.MessageVariableCompile,
					},
				},
			},
		},
		{
			name: "ConstraintError_Compile",
			in: domain.VariableProcessIn{
				Variables: []domain.Variable{
					{
						ID:      1,
						Name:    "var1",
						Type:    variable_domain.TypeInteger,
						IsInput: true,
						Constraints: []domain.Constraint{
							{
								ID:         1,
								Name:       "expr1",
								Expression: "invalid",
								IsActive:   true,
							},
							{
								ID:         2,
								Name:       "expr2",
								Expression: "invalid",
								IsActive:   true,
							},
						},
					},
				},
				Payload: map[string]string{"var1": "100"},
			},
			want: task_domain.ProcessError{
				VariableErrors: []task_domain.VariableError{
					{
						ID:   1,
						Name: "var1",
						ConstraintErrors: []task_domain.ConstraintError{
							{
								ID:      1,
								Name:    "expr1",
								Message: task_domain.MessageConstraintCompile,
							},
							{
								ID:      2,
								Name:    "expr2",
								Message: task_domain.MessageConstraintCompile,
							},
						},
					},
				},
			},
		},
		{
			name: "ConstraintError_Check",
			in: domain.VariableProcessIn{
				Variables: []domain.Variable{
					{
						ID:      1,
						Name:    "var1",
						Type:    variable_domain.TypeInteger,
						IsInput: true,
						Constraints: []domain.Constraint{
							{
								ID:         1,
								Name:       "expr1",
								Expression: "var1 > 200",
								IsActive:   true,
							},
							{
								ID:         2,
								Name:       "expr2",
								Expression: "var1 > 300",
								IsActive:   false,
							},
							{
								ID:         3,
								Name:       "expr3",
								Expression: "var1 > 400",
								IsActive:   true,
							},
						},
					},
				},
				Payload: map[string]string{"var1": "100"},
			},
			want: task_domain.ProcessError{
				VariableErrors: []task_domain.VariableError{
					{
						ID:   1,
						Name: "var1",
						ConstraintErrors: []task_domain.ConstraintError{
							{
								ID:      1,
								Name:    "expr1",
								Message: task_domain.MessageConstraintCheck,
							},
							{
								ID:      3,
								Name:    "expr3",
								Message: task_domain.MessageConstraintCheck,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Handle(ctx, tt.in)

			var got *task_domain.ProcessError
			require.ErrorAs(t, err, &got)
			require.Equal(t, tt.want, *got)
		})
	}
}
