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
