package variable_process_service

import (
	"context"
	"math"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

func TestBuiltins_Math(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       any
	}{
		{"sqrt", "sqrt(9)", 3.0},
		{"pow", "pow(2, 10)", 1024.0},
		{"exp_log", "round(log(exp(1)), 4)", 1.0},
		{"log10", "log10(1000)", 3.0},
		{"log2", "log2(8)", 3.0},
		{"sin_pi_half", "round(sin(pi / 2), 4)", 1.0},
		{"cos_pi", "round(cos(pi), 4)", -1.0},
		{"tan_zero", "round(tan(0), 4)", 0.0},
		{"asin_one", "round(asin(1) - pi / 2, 4)", 0.0},
		{"acos_one", "round(acos(1), 4)", 0.0},
		{"atan_one", "round(atan(1) - pi / 4, 4)", 0.0},
		{"atan2", "round(atan2(1, 1) - pi / 4, 4)", 0.0},
		{"hypot", "hypot(3, 4)", 5.0},
		{"mod", "mod(10, 3)", 1.0},
		{"sign_pos", "sign(42)", 1.0},
		{"sign_neg", "sign(-3.14)", -1.0},
		{"sign_zero", "sign(0)", 0.0},
		{"clamp_low", "clamp(-5, 0, 10)", 0.0},
		{"clamp_high", "clamp(50, 0, 10)", 10.0},
		{"clamp_mid", "clamp(5, 0, 10)", 5.0},
		{"interpolate", "interpolate(2.5, 0, 5, 0, 100)", 50.0},
		{"const_e", "round(e, 4)", 2.7183},
		{"const_g", "g", 9.80665},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runExpression(t, tt.expression, variable_domain.TypeFloat)
			require.InDelta(t, tt.want, got, 1e-6)
		})
	}
}

func TestBuiltins_RoundAndFormat(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		typ        variable_domain.Type
		want       any
	}{
		{"round_no_arg", "round(3.7)", variable_domain.TypeFloat, 4.0},
		{"round_two_decimals", "round(3.14159, 2)", variable_domain.TypeFloat, 3.14},
		{"round_zero_decimals", "round(3.7, 0)", variable_domain.TypeFloat, 4.0},
		{"round_negative", "round(-3.55, 1)", variable_domain.TypeFloat, -3.6},
		{"round_step_half", "roundStep(12.3, 0.5)", variable_domain.TypeFloat, 12.5},
		{"round_step_hundred", "roundStep(1234, 100)", variable_domain.TypeFloat, 1200.0},
		{"format_number_default", `formatNumber(1234567.5, 2)`, variable_domain.TypeString, "1 234 567,50"},
		{"format_number_custom", `formatNumber(1234.5, 1, ".", ",")`, variable_domain.TypeString, "1,234.5"},
		{"format_number_no_thousands", `formatNumber(1234.5, 2, ",", "")`, variable_domain.TypeString, "1234,50"},
		{"format_number_negative", `formatNumber(-1234.5, 1)`, variable_domain.TypeString, "-1 234,5"},
		{"percent_zero", `percent(0.157)`, variable_domain.TypeString, "16%"},
		{"percent_one_decimal", `percent(0.157, 1)`, variable_domain.TypeString, "15,7%"},
		{"scientific_default", `scientific(123000.0)`, variable_domain.TypeString, "1.23e+05"},
		{"scientific_zero", `scientific(123000.0, 0)`, variable_domain.TypeString, "1e+05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runExpression(t, tt.expression, tt.typ)
			switch want := tt.want.(type) {
			case float64:
				require.InDelta(t, want, got, 1e-9)
			default:
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestBuiltins_ConstantsStrippedFromOutput(t *testing.T) {
	ctx := context.Background()
	service := New()

	in := domain.VariableProcessIn{
		Variables: []domain.Variable{
			{
				ID:         1,
				Name:       "circle_area",
				Type:       variable_domain.TypeFloat,
				Expression: lo.ToPtr("pi * pow(radius, 2)"),
			},
			{
				ID:      2,
				Name:    "radius",
				Type:    variable_domain.TypeFloat,
				IsInput: true,
			},
		},
		Payload: map[string]string{"radius": "2"},
	}

	got, err := service.Handle(ctx, in)
	require.NoError(t, err)

	require.NotContains(t, got, "pi")
	require.NotContains(t, got, "e")
	require.NotContains(t, got, "g")

	require.InDelta(t, 4*math.Pi, got["circle_area"], 1e-9)
}

func TestBuiltins_ConstraintsSeeConstants(t *testing.T) {
	ctx := context.Background()
	service := New()

	in := domain.VariableProcessIn{
		Variables: []domain.Variable{
			{
				ID:      1,
				Name:    "angle",
				Type:    variable_domain.TypeFloat,
				IsInput: true,
				Constraints: []domain.Constraint{
					{
						ID:         1,
						Name:       "below_two_pi",
						Expression: "angle < 2 * pi",
						IsActive:   true,
					},
				},
			},
		},
		Payload: map[string]string{"angle": "5"},
	}

	_, err := service.Handle(ctx, in)
	require.NoError(t, err)
}

func TestToFloat_AllTypes(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want float64
	}{
		{"float64", float64(1.5), 1.5},
		{"float32", float32(1.5), 1.5},
		{"int", int(42), 42},
		{"int8", int8(-5), -5},
		{"int16", int16(-300), -300},
		{"int32", int32(-70000), -70000},
		{"int64", int64(1 << 40), 1 << 40},
		{"uint", uint(42), 42},
		{"uint8", uint8(255), 255},
		{"uint16", uint16(65535), 65535},
		{"uint32", uint32(1 << 30), 1 << 30},
		{"uint64", uint64(1 << 40), 1 << 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toFloat(tt.in)
			require.NoError(t, err)
			require.InDelta(t, tt.want, got, 1e-9)
		})
	}
}

func TestToFloat_Unsupported(t *testing.T) {
	tests := []any{"string", true, nil, struct{}{}, []int{1, 2}}

	for _, in := range tests {
		_, err := toFloat(in)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot convert")
	}
}

func TestToInt_AllTypes(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int
	}{
		{"int", int(42), 42},
		{"int8", int8(-5), -5},
		{"int16", int16(-300), -300},
		{"int32", int32(-70000), -70000},
		{"int64", int64(1 << 30), 1 << 30},
		{"uint8", uint8(255), 255},
		{"uint16", uint16(65535), 65535},
		{"uint32", uint32(1 << 30), 1 << 30},
		{"float32", float32(3.7), 3},
		{"float64", float64(-3.7), -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToInt_Unsupported(t *testing.T) {
	tests := []any{"string", true, nil, uint(1), uint64(1)}

	for _, in := range tests {
		_, err := toInt(in)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot convert")
	}
}

func TestFn1_TypeError(t *testing.T) {
	wrapped := fn1(math.Sqrt)
	_, err := wrapped("not a number")
	require.Error(t, err)
}

func TestFn2_TypeError(t *testing.T) {
	wrapped := fn2(math.Pow)

	_, err := wrapped("nope", 2.0)
	require.Error(t, err)

	_, err = wrapped(2.0, "nope")
	require.Error(t, err)
}

func TestFuncRound_Errors(t *testing.T) {
	_, err := funcRound("not a number")
	require.Error(t, err)

	_, err = funcRound(1.5, "not an int")
	require.Error(t, err)
}

func TestFuncRound_NegativeDigits(t *testing.T) {
	got, err := funcRound(1234.0, -2)
	require.NoError(t, err)
	require.InDelta(t, 1200.0, got, 1e-9)
}

func TestFuncRoundStep_Errors(t *testing.T) {
	_, err := funcRoundStep("nope", 1.0)
	require.Error(t, err)

	_, err = funcRoundStep(1.0, "nope")
	require.Error(t, err)

	_, err = funcRoundStep(1.0, 0.0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "step must be non-zero")
}

func TestFuncClamp_Errors(t *testing.T) {
	_, err := funcClamp("x", 0.0, 10.0)
	require.Error(t, err)

	_, err = funcClamp(1.0, "x", 10.0)
	require.Error(t, err)

	_, err = funcClamp(1.0, 0.0, "x")
	require.Error(t, err)
}

func TestFuncInterpolate_Errors(t *testing.T) {
	_, err := funcInterpolate(1.0, 0.0, 1.0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "expected 5 arguments")

	_, err = funcInterpolate("nope", 0.0, 1.0, 0.0, 100.0)
	require.Error(t, err)

	_, err = funcInterpolate(0.5, 1.0, 1.0, 0.0, 100.0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "x0 and x1 must differ")
}

func TestFuncFormatNumber_Errors(t *testing.T) {
	_, err := funcFormatNumber("x", 2)
	require.Error(t, err)

	_, err = funcFormatNumber(1.0, "x")
	require.Error(t, err)

	_, err = funcFormatNumber(1.0, 2, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "decimal separator must be a string")

	_, err = funcFormatNumber(1.0, 2, ",", 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "thousand separator must be a string")
}

func TestFormatNumber_Edge(t *testing.T) {
	require.Equal(t, "0", formatNumber(0, -1, ",", " "))
	require.Equal(t, "5", formatNumber(5, 0, ",", " "))
	require.Equal(t, "1,5", formatNumber(1.5, 1, ",", " "))
	require.Equal(t, "100", formatNumber(100, 0, ",", " "))
	require.Equal(t, "999", formatNumber(999, 0, ",", " "))
}

func TestFuncPercent_Errors(t *testing.T) {
	_, err := funcPercent("nope")
	require.Error(t, err)

	_, err = funcPercent(0.5, "nope")
	require.Error(t, err)
}

func TestFuncScientific_Errors(t *testing.T) {
	_, err := funcScientific("nope")
	require.Error(t, err)

	_, err = funcScientific(1.0, "nope")
	require.Error(t, err)
}

func TestFuncScientific_NegativeDecimals(t *testing.T) {
	got, err := funcScientific(123000.0, -3)
	require.NoError(t, err)
	require.Equal(t, "1e+05", got)
}

// runExpression evaluates a single expression by wrapping it into a computed
// variable so the existing Service pipeline does the heavy lifting.
func runExpression(t *testing.T, expression string, typ variable_domain.Type) any {
	t.Helper()
	ctx := context.Background()
	service := New()

	in := domain.VariableProcessIn{
		Variables: []domain.Variable{
			{
				ID:         1,
				Name:       "result",
				Type:       typ,
				Expression: lo.ToPtr(expression),
			},
		},
	}

	got, err := service.Handle(ctx, in)
	if err != nil {
		var processErr *task_domain.ProcessError
		require.ErrorAs(t, err, &processErr)
		t.Fatalf("expression %q failed: %+v", expression, *processErr)
	}

	return got["result"]
}
