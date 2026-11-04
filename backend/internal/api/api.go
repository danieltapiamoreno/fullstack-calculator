// Package api exposes the calculator over HTTP as a JSON REST API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"fullstack-calculator/backend/internal/calculator"
)

const (
	routePrefix  = "/api/v1/"
	maxBodyBytes = 1024

	codeInvalidRequest   = "invalid_request"
	codeRequestTooLarge  = "request_too_large"
	codeDivisionByZero   = "division_by_zero"
	codeUndefinedResult  = "undefined_result"
	codeResultOutOfRange = "result_out_of_range"
	codeInternalError    = "internal_error"
)

// operations maps each route name to its handler.
var operations = map[string]http.HandlerFunc{
	"add":        binary(calculator.Add),
	"subtract":   binary(calculator.Subtract),
	"multiply":   binary(calculator.Multiply),
	"divide":     binary(calculator.Divide),
	"power":      binary(calculator.Power),
	"sqrt":       unary(calculator.Sqrt),
	"percentage": binary(calculator.Percentage),
}

// NewHandler returns a handler serving every calculator route and the health
// check, plus the files in staticDir unless it is empty.
func NewHandler(staticDir string) http.Handler {
	mux := http.NewServeMux()
	for name, handler := range operations {
		mux.Handle("POST "+routePrefix+name, handler)
	}
	mux.HandleFunc("GET /healthz", health)
	if staticDir == "" {
		return mux
	}

	// A "GET /" pattern on mux would match every API path too, turning its
	// default 404 and 405 responses into file responses, so the API's paths are
	// sent to mux before anything can reach the files.
	withFiles := http.NewServeMux()
	withFiles.Handle("/api", mux)
	withFiles.Handle("/api/", mux)
	withFiles.Handle("/healthz", mux)
	withFiles.Handle("/", staticFiles(staticDir))
	return withFiles
}

// staticFiles serves the files in dir to GET and HEAD requests only, because
// http.FileServer alone would answer any method. It has its own mux because
// "GET /" conflicts with the method-less API patterns it would otherwise share
// one with.
func staticFiles(dir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(dir)))
	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{"ok"})
}

// binary adapts a two-operand calculation into an HTTP handler.
func binary(calculate func(a, b float64) (float64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operands, ok := decodeOperands(w, r, "a", "b")
		if !ok {
			return
		}
		result, err := calculate(operands["a"], operands["b"])
		respond(w, result, err)
	}
}

// unary adapts a one-operand calculation into an HTTP handler.
func unary(calculate func(a float64) (float64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operands, ok := decodeOperands(w, r, "a")
		if !ok {
			return
		}
		result, err := calculate(operands["a"])
		respond(w, result, err)
	}
}

// decodeOperands reads a JSON object containing exactly the named numeric
// fields, writing an error response and returning false if the body is anything
// else. It reads the whole body first so that the size limit applies however the
// body is malformed, and decodes into a map because map keys match exactly,
// whereas struct fields ignore case. Pointer values tell a null from a zero.
func decodeOperands(w http.ResponseWriter, r *http.Request, names ...string) (map[string]float64, bool) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, codeRequestTooLarge, "the request body is too large")
			return nil, false
		}
		writeInvalidRequest(w, names)
		return nil, false
	}

	var fields map[string]*float64
	if err := json.Unmarshal(body, &fields); err != nil || len(fields) != len(names) {
		writeInvalidRequest(w, names)
		return nil, false
	}

	operands := make(map[string]float64, len(names))
	for _, name := range names {
		value := fields[name]
		if value == nil {
			writeInvalidRequest(w, names)
			return nil, false
		}
		operands[name] = *value
	}
	return operands, true
}

func writeInvalidRequest(w http.ResponseWriter, names []string) {
	message := "the request body must be a JSON object containing only the numeric fields: " + strings.Join(names, ", ")
	writeError(w, http.StatusBadRequest, codeInvalidRequest, message)
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
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, codeDivisionByZero, err.Error())
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
