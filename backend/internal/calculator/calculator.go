// Package calculator implements the calculator's arithmetic. It knows nothing
// about HTTP or JSON, and every operation returns either a finite result or one
// of the errors defined here.
package calculator

import (
	"errors"
	"math"
)

var (
	// ErrDivisionByZero means the operation would divide by zero.
	ErrDivisionByZero = errors.New("cannot divide by zero")
	// ErrUndefinedResult means the operation has no real-number answer.
	ErrUndefinedResult = errors.New("the result is not a real number")
	// ErrResultOutOfRange means the answer is too large to represent.
	ErrResultOutOfRange = errors.New("the result is too large to represent")
)

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	return checked(a + b)
}

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	return checked(a - b)
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	return checked(a * b)
}

// Divide returns a / b. Any zero b is an error, including when a is also zero,
// and -0 compares equal to 0.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return checked(a / b)
}

// Percentage returns a percent of b. It calculates a * b / 100 on the operands'
// fractions and restores their exponents last, so the product can't overflow
// on its own, as a * b would on 1e308 % of 50. Only a result within rounding
// error of the largest float64 can overflow when its true value fits. Multiplying
// first keeps common results exact, where a / 100 * b makes 7 % of 100 into
// 7.000000000000001.
func Percentage(a, b float64) (float64, error) {
	aFraction, aExponent := math.Frexp(a)
	bFraction, bExponent := math.Frexp(b)
	return checked(math.Ldexp(aFraction*bFraction/100, aExponent+bExponent))
}

// checked rejects non-finite results and normalises -0 to 0. Every operation
// returns through it, so callers never receive a value JSON can't represent.
func checked(result float64) (float64, error) {
	switch {
	case math.IsNaN(result):
		return 0, ErrUndefinedResult
	case math.IsInf(result, 0):
		return 0, ErrResultOutOfRange
	case result == 0:
		return 0, nil
	default:
		return result, nil
	}
}
