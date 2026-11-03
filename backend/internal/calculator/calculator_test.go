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
