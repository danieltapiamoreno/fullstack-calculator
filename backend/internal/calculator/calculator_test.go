package calculator_test

import (
	"errors"
	"math"
	"testing"

	"fullstack-calculator/backend/internal/calculator"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "whole numbers", a: 2, b: 3, want: 5},
		{name: "negative numbers", a: -2, b: -3, want: -5},
		{name: "decimals", a: 1.5, b: 2.25, want: 3.75},
		{name: "keeps floating-point precision", a: 0.1, b: 0.2, want: 0.30000000000000004},
		{name: "overflow", a: 1e308, b: 1e308, wantErr: calculator.ErrResultOutOfRange},
		{name: "negative overflow", a: -1e308, b: -1e308, wantErr: calculator.ErrResultOutOfRange},
		{name: "not a number", a: math.Inf(1), b: math.Inf(-1), wantErr: calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Add(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Add(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAddNeverReturnsNegativeZero(t *testing.T) {
	negativeZero := math.Copysign(0, -1)

	got, err := calculator.Add(negativeZero, negativeZero)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Signbit(got) {
		t.Errorf("Add(-0, -0) = %v, want positive zero", got)
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "whole numbers", a: 2, b: 3, want: -1},
		{name: "subtracting a negative", a: 2, b: -3, want: 5},
		{name: "decimals", a: 5.5, b: 2.25, want: 3.25},
		{name: "keeps floating-point precision", a: 0.3, b: 0.1, want: 0.19999999999999998},
		{name: "overflow", a: 1e308, b: -1e308, wantErr: calculator.ErrResultOutOfRange},
		{name: "negative overflow", a: -1e308, b: 1e308, wantErr: calculator.ErrResultOutOfRange},
		{name: "not a number", a: math.Inf(1), b: math.Inf(1), wantErr: calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Subtract(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Subtract(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtractNeverReturnsNegativeZero(t *testing.T) {
	got, err := calculator.Subtract(math.Copysign(0, -1), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Signbit(got) {
		t.Errorf("Subtract(-0, 0) = %v, want positive zero", got)
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "whole numbers", a: 2, b: 3, want: 6},
		{name: "negative numbers", a: -2, b: 3, want: -6},
		{name: "decimals", a: 1.5, b: 2.5, want: 3.75},
		{name: "keeps floating-point precision", a: 0.1, b: 0.2, want: 0.020000000000000004},
		{name: "overflow", a: 1e308, b: 10, wantErr: calculator.ErrResultOutOfRange},
		{name: "negative overflow", a: 1e308, b: -10, wantErr: calculator.ErrResultOutOfRange},
		{name: "not a number", a: 0, b: math.Inf(1), wantErr: calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Multiply(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Multiply(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiplyNeverReturnsNegativeZero(t *testing.T) {
	got, err := calculator.Multiply(0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Signbit(got) {
		t.Errorf("Multiply(0, -1) = %v, want positive zero", got)
	}
}

func TestDivide(t *testing.T) {
	negativeZero := math.Copysign(0, -1)

	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "whole numbers", a: 6, b: 3, want: 2},
		{name: "fractional result", a: 7, b: 2, want: 3.5},
		{name: "negative dividend", a: -6, b: 3, want: -2},
		{name: "keeps floating-point precision", a: 1, b: 3, want: 0.3333333333333333},
		{name: "zero dividend", a: 0, b: 5, want: 0},
		{name: "zero divisor", a: 1, b: 0, wantErr: calculator.ErrDivisionByZero},
		{name: "negative dividend and zero divisor", a: -1, b: 0, wantErr: calculator.ErrDivisionByZero},
		{name: "zero divided by zero", a: 0, b: 0, wantErr: calculator.ErrDivisionByZero},
		{name: "negative zero divisor", a: 1, b: negativeZero, wantErr: calculator.ErrDivisionByZero},
		{name: "zero divided by negative zero", a: 0, b: negativeZero, wantErr: calculator.ErrDivisionByZero},
		{name: "overflow", a: 1e308, b: 1e-308, wantErr: calculator.ErrResultOutOfRange},
		{name: "negative overflow", a: -1e308, b: 1e-308, wantErr: calculator.ErrResultOutOfRange},
		{name: "not a number", a: math.Inf(1), b: math.Inf(1), wantErr: calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Divide(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Divide(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivideNeverReturnsNegativeZero(t *testing.T) {
	got, err := calculator.Divide(0, -5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Signbit(got) {
		t.Errorf("Divide(0, -5) = %v, want positive zero", got)
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "whole numbers", a: 20, b: 50, want: 10},
		{name: "whole percent of 100", a: 7, b: 100, want: 7},
		{name: "another whole percent of 100", a: 29, b: 100, want: 29},
		{name: "result with no exact binary form", a: 7, b: 3, want: 0.21},
		{name: "decimal percent", a: 12.5, b: 80, want: 10},
		{name: "decimal amount keeps floating-point precision", a: 15, b: 19.99, want: 2.9984999999999995},
		{name: "both decimals keep floating-point precision", a: 0.1, b: 0.2, want: 0.00020000000000000004},
		{name: "negative percent", a: -20, b: 50, want: -10},
		{name: "negative amount", a: 20, b: -50, want: -10},
		{name: "both negative", a: -20, b: -50, want: 10},
		{name: "percent over 100", a: 250, b: 40, want: 100},
		{name: "zero percent", a: 0, b: 50, want: 0},
		{name: "percent of zero", a: 20, b: 0, want: 0},
		{name: "large percent whose product overflows", a: 1e308, b: 50, want: 5e307},
		{name: "large amount whose product overflows", a: 50, b: 1e308, want: 5e307},
		{name: "largest result", a: math.MaxFloat64, b: 100, want: math.MaxFloat64},
		{name: "largest negative result", a: -math.MaxFloat64, b: 100, want: -math.MaxFloat64},
		{name: "smallest value is kept", a: 5e-324, b: 100, want: 5e-324},
		{name: "underflows to zero", a: 1e-200, b: 1e-200, want: 0},
		{name: "overflow", a: 1e308, b: 1000, wantErr: calculator.ErrResultOutOfRange},
		{name: "overflow just past the largest result", a: math.MaxFloat64, b: 101, wantErr: calculator.ErrResultOutOfRange},
		{name: "negative overflow", a: -1e308, b: 1e308, wantErr: calculator.ErrResultOutOfRange},
		{name: "infinite input", a: math.Inf(1), b: 50, wantErr: calculator.ErrResultOutOfRange},
		{name: "zero percent of infinity", a: 0, b: math.Inf(1), wantErr: calculator.ErrUndefinedResult},
		{name: "not a number", a: math.NaN(), b: 50, wantErr: calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Percentage(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Percentage(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPercentageNeverReturnsNegativeZero(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
	}{
		{name: "zero percent of a negative", a: 0, b: -50},
		{name: "negative result that underflows", a: -1e-200, b: 1e-200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Percentage(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Signbit(got) {
				t.Errorf("Percentage(%v, %v) = %v, want positive zero", tt.a, tt.b, got)
			}
		})
	}
}
