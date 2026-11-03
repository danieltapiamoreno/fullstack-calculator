// Package api exposes the calculator over HTTP as a JSON REST API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"fullstack-calculator/backend/internal/calculator"
)

const (
	routePrefix  = "/api/v1/"
	maxBodyBytes = 1024

	codeInvalidRequest   = "invalid_request"
	codeRequestTooLarge  = "request_too_large"
	codeUndefinedResult  = "undefined_result"
	codeResultOutOfRange = "result_out_of_range"
	codeInternalError    = "internal_error"
)

// operations maps each route name to its handler.
var operations = map[string]http.HandlerFunc{
	"add": binary(calculator.Add),
}

// NewHandler returns a handler serving every calculator route.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	for name, handler := range operations {
		mux.Handle("POST "+routePrefix+name, handler)
	}
	return mux
}

// binary adapts a two-operand calculation into an HTTP handler.
func binary(calculate func(a, b float64) (float64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			A *float64 `json:"a"`
			B *float64 `json:"b"`
		}
		if !decode(w, r, &request) {
			return
		}
		if request.A == nil || request.B == nil {
			writeError(w, http.StatusBadRequest, codeInvalidRequest, `the fields "a" and "b" are required`)
			return
		}
		result, err := calculate(*request.A, *request.B)
		respond(w, result, err)
	}
}

// decode reads the request body into dst, writing an error response and
// returning false if the body isn't a single valid JSON value of dst's shape.
// Unset fields stay nil, which lets handlers tell a missing value from zero.
func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err == nil {
		err = expectEOF(decoder)
	}
	if err == nil {
		return true
	}

	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, codeRequestTooLarge, "the request body is too large")
	} else {
		writeError(w, http.StatusBadRequest, codeInvalidRequest, "the request body must be a JSON object containing only the expected numeric fields")
	}
	return false
}

// expectEOF reports an error if anything other than whitespace follows the
// value the decoder has just read.
func expectEOF(decoder *json.Decoder) error {
	err := decoder.Decode(&struct{}{})
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err == nil {
		return errors.New("unexpected data after JSON value")
	}
	return err
}

// respond writes a calculation's outcome.
func respond(w http.ResponseWriter, result float64, err error) {
	if err != nil {
		writeCalculationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Result float64 `json:"result"`
	}{result})
}

func writeCalculationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrUndefinedResult):
		writeError(w, http.StatusUnprocessableEntity, codeUndefinedResult, err.Error())
	case errors.Is(err, calculator.ErrResultOutOfRange):
		writeError(w, http.StatusUnprocessableEntity, codeResultOutOfRange, err.Error())
	default:
		log.Printf("unexpected calculation error: %v", err)
		writeError(w, http.StatusInternalServerError, codeInternalError, "something went wrong")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	type detail struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	writeJSON(w, status, struct {
		Error detail `json:"error"`
	}{detail{code, message}})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writing response: %v", err)
	}
}
