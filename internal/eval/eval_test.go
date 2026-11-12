package eval_test

import (
	"math"
	"strings"
	"testing"

	"github.com/gleb/go-calc/internal/eval"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{name: "zero", expr: "0", want: 0},
		{name: "simple addition", expr: "2 + 2", want: 4},
		{name: "multiplication priority", expr: "2 + 2 * 2", want: 6},
		{name: "subtraction and division", expr: "10 - 4 / 2", want: 8},
		{name: "modulo", expr: "10 % 3", want: 1},
		{name: "unary minus", expr: "-5", want: -5},
		{name: "double unary minus", expr: "--5", want: 5},
		{name: "subtraction of negative", expr: "2--2", want: 4},
		{name: "parentheses", expr: "(2 + 3) * 4", want: 20},
		{name: "unary minus with parens", expr: "-(2+3)*4", want: -20},
		{name: "power right associativity", expr: "2 ^ 3 ^ 2", want: 512}, // 2^(3^2) = 2^9 = 512
		{name: "negative power", expr: "2^-1", want: 0.5},
		{name: "unary minus before power", expr: "-2^2", want: -4},
		{name: "parens negative base power", expr: "(-2)^2", want: 4},
		{name: "scientific notation", expr: "1e-3 + 2", want: 2.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.Eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestEval_Errors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{name: "empty string", expr: "", wantErr: "empty expression"},
		{name: "division by zero", expr: "1 / 0", wantErr: "division by zero"},
		{name: "zero to negative power", expr: "0 ^ -1", wantErr: "division by zero"},
		{name: "modulo by zero", expr: "5 % 0", wantErr: "modulo by zero"},
		{name: "overflow", expr: "1e300 * 1e300", wantErr: "infinity"},
		{name: "NaN", expr: "(-2) ^ 0.5", wantErr: "NaN"},
		{name: "unclosed parens", expr: "(2 + 3", wantErr: "missing closing parenthesis"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := eval.Eval(tt.expr)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tt.expr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		val       float64
		precision int
		want      string
	}{
		{val: 0, precision: 10, want: "0"},
		{val: 4.0, precision: 10, want: "4"},
		{val: 4.5, precision: 10, want: "4.5"},
		{val: 1.0 / 3.0, precision: 2, want: "0.33"},
		{val: 1.0 / 3.0, precision: 4, want: "0.3333"},
		{val: 1e-3, precision: 10, want: "0.001"},
		{val: 1e18, precision: 4, want: "1e+18"},
	}

	for _, tt := range tests {
		got := eval.Format(tt.val, tt.precision)
		if got != tt.want {
			t.Errorf("Format(%v, %d) = %q, want %q", tt.val, tt.precision, got, tt.want)
		}
	}
}
