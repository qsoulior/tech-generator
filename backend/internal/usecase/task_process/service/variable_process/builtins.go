package variable_process_service

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
)

// mathConstants are injected into the evaluation environment so users can write
// "2 * pi * r" instead of "pi() * ...". They are stripped from the returned
// variable map so they never leak into the rendered template data.
var mathConstants = map[string]any{
	"pi": math.Pi,
	"e":  math.E,
	"g":  9.80665,
}

var builtinOptions = []expr.Option{
	expr.Function("sqrt", fn1(math.Sqrt)),
	expr.Function("exp", fn1(math.Exp)),
	expr.Function("log", fn1(math.Log)),
	expr.Function("log10", fn1(math.Log10)),
	expr.Function("log2", fn1(math.Log2)),
	expr.Function("sin", fn1(math.Sin)),
	expr.Function("cos", fn1(math.Cos)),
	expr.Function("tan", fn1(math.Tan)),
	expr.Function("asin", fn1(math.Asin)),
	expr.Function("acos", fn1(math.Acos)),
	expr.Function("atan", fn1(math.Atan)),
	expr.Function("sign", fn1(funcSign)),

	expr.Function("pow", fn2(math.Pow)),
	expr.Function("atan2", fn2(math.Atan2)),
	expr.Function("hypot", fn2(math.Hypot)),
	expr.Function("mod", fn2(math.Mod)),

	expr.Function("clamp", funcClamp),
	expr.Function("interpolate", funcInterpolate),
	expr.Function("round", funcRound),
	expr.Function("roundStep", funcRoundStep),
	expr.Function("formatNumber", funcFormatNumber),
	expr.Function("percent", funcPercent),
	expr.Function("scientific", funcScientific),
}

func fn1(f func(float64) float64) func(...any) (any, error) {
	return func(params ...any) (any, error) {
		x, err := toFloat(params[0])
		if err != nil {
			return nil, err
		}
		return f(x), nil
	}
}

func fn2(f func(float64, float64) float64) func(...any) (any, error) {
	return func(params ...any) (any, error) {
		x, err := toFloat(params[0])
		if err != nil {
			return nil, err
		}
		y, err := toFloat(params[1])
		if err != nil {
			return nil, err
		}
		return f(x, y), nil
	}
}

func funcSign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func funcClamp(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	low, err := toFloat(params[1])
	if err != nil {
		return nil, err
	}
	high, err := toFloat(params[2])
	if err != nil {
		return nil, err
	}
	return math.Max(low, math.Min(high, x)), nil
}

func funcInterpolate(params ...any) (any, error) {
	xs := make([]float64, len(params))
	for i, p := range params {
		v, err := toFloat(p)
		if err != nil {
			return nil, err
		}
		xs[i] = v
	}
	if len(xs) != 5 {
		return nil, fmt.Errorf("interpolate: expected 5 arguments, got %d", len(xs))
	}
	x, x0, x1, y0, y1 := xs[0], xs[1], xs[2], xs[3], xs[4]
	if x1 == x0 {
		return nil, fmt.Errorf("interpolate: x0 and x1 must differ")
	}
	return y0 + (x-x0)*(y1-y0)/(x1-x0), nil
}

func funcRound(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	if len(params) < 2 {
		return math.Round(x), nil
	}
	n, err := toInt(params[1])
	if err != nil {
		return nil, err
	}
	m := math.Pow(10, float64(n))
	return math.Round(x*m) / m, nil
}

func funcRoundStep(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	step, err := toFloat(params[1])
	if err != nil {
		return nil, err
	}
	if step == 0 {
		return nil, fmt.Errorf("roundStep: step must be non-zero")
	}
	return math.Round(x/step) * step, nil
}

func funcFormatNumber(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	decimals, err := toInt(params[1])
	if err != nil {
		return nil, err
	}
	// Default decimal separator is comma and thousand separator is U+00A0
	// (NBSP) so digit groups don't wrap in rendered docs — matches Russian
	// typographic convention.
	decSep := ","
	thouSep := " "
	if len(params) >= 3 {
		s, ok := params[2].(string)
		if !ok {
			return nil, fmt.Errorf("formatNumber: decimal separator must be a string")
		}
		decSep = s
	}
	if len(params) >= 4 {
		s, ok := params[3].(string)
		if !ok {
			return nil, fmt.Errorf("formatNumber: thousand separator must be a string")
		}
		thouSep = s
	}
	return formatNumber(x, decimals, decSep, thouSep), nil
}

func formatNumber(x float64, decimals int, decSep, thouSep string) string {
	if decimals < 0 {
		decimals = 0
	}
	raw := strconv.FormatFloat(x, 'f', decimals, 64)
	sign := ""
	if strings.HasPrefix(raw, "-") {
		sign = "-"
		raw = raw[1:]
	}
	intPart, fracPart, _ := strings.Cut(raw, ".")
	if thouSep != "" && len(intPart) > 3 {
		var b strings.Builder
		first := len(intPart) % 3
		if first > 0 {
			b.WriteString(intPart[:first])
		}
		for i := first; i < len(intPart); i += 3 {
			if b.Len() > 0 {
				b.WriteString(thouSep)
			}
			b.WriteString(intPart[i : i+3])
		}
		intPart = b.String()
	}
	if fracPart != "" {
		return sign + intPart + decSep + fracPart
	}
	return sign + intPart
}

func funcPercent(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	decimals := 0
	if len(params) >= 2 {
		decimals, err = toInt(params[1])
		if err != nil {
			return nil, err
		}
	}
	return formatNumber(x*100, decimals, ",", "") + "%", nil
}

func funcScientific(params ...any) (any, error) {
	x, err := toFloat(params[0])
	if err != nil {
		return nil, err
	}
	decimals := 2
	if len(params) >= 2 {
		decimals, err = toInt(params[1])
		if err != nil {
			return nil, err
		}
	}
	if decimals < 0 {
		decimals = 0
	}
	return strconv.FormatFloat(x, 'e', decimals, 64), nil
}

func toFloat(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	}
	return 0, fmt.Errorf("cannot convert %T to float64", v)
}

func toInt(v any) (int, error) {
	switch x := v.(type) {
	case int:
		return x, nil
	case int8:
		return int(x), nil
	case int16:
		return int(x), nil
	case int32:
		return int(x), nil
	case int64:
		return int(x), nil
	case uint8:
		return int(x), nil
	case uint16:
		return int(x), nil
	case uint32:
		return int(x), nil
	case float32:
		return int(x), nil
	case float64:
		return int(x), nil
	}
	return 0, fmt.Errorf("cannot convert %T to int", v)
}
