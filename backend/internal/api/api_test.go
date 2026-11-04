package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fullstack-calculator/backend/internal/api"
)

type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func send(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	api.NewHandler().ServeHTTP(recorder, request)
	return recorder
}

func TestAddReturnsResult(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/add", `{"a": 2, "b": 3}`)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := strings.TrimSpace(response.Body.String()); got != `{"result":5}` {
		t.Errorf("body = %s, want {\"result\":5}", got)
	}
}

func TestRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"a": 2,`},
		{name: "empty body", body: ``},
		{name: "JSON null", body: `null`},
		{name: "not an object", body: `[2, 3]`},
		{name: "missing a", body: `{"b": 3}`},
		{name: "missing b", body: `{"a": 2}`},
		{name: "null field", body: `{"a": 2, "b": null}`},
		{name: "string field", body: `{"a": "2", "b": 3}`},
		{name: "unknown field", body: `{"a": 2, "b": 3, "c": 4}`},
		{name: "number too large for float64", body: `{"a": 1e999, "b": 3}`},
		{name: "trailing data", body: `{"a": 2, "b": 3} {"a": 4}`},
		{name: "field names in a different case", body: `{"A": 2, "B": 3}`},
		{name: "extra field differing only in case", body: `{"a": 2, "b": 3, "A": 10}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/add", tt.body)
			assertError(t, response, http.StatusBadRequest, "invalid_request")
		})
	}
}

func TestRejectsOversizedBodies(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "oversized value", body: `{"a": 1` + strings.Repeat("0", 1024) + `, "b": 3}`},
		{name: "valid request followed by padding", body: `{"a": 2, "b": 3}` + strings.Repeat(" ", 1024)},
		{name: "wrong type followed by padding", body: `{"a": "2", "b": 3}` + strings.Repeat(" ", 2048)},
		{name: "unknown field followed by padding", body: `{"a": 2, "b": 3, "c": 4}` + strings.Repeat(" ", 2048)},
		{name: "malformed JSON followed by padding", body: `x` + strings.Repeat(" ", 2048)},
		{name: "valid request followed by data and padding", body: `{"a": 2, "b": 3}x` + strings.Repeat(" ", 2048)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/add", tt.body)
			assertError(t, response, http.StatusRequestEntityTooLarge, "request_too_large")
		})
	}
}

func TestAcceptsBodyAtSizeLimit(t *testing.T) {
	body := `{"a": 2, "b": 3}`
	body += strings.Repeat(" ", 1024-len(body))

	response := send(t, http.MethodPost, "/api/v1/add", body)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
}

func TestOverflowReturnsResultOutOfRange(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/add", `{"a": 1e308, "b": 1e308}`)

	assertError(t, response, http.StatusUnprocessableEntity, "result_out_of_range")
}

func TestNegativeZeroIsReturnedAsZero(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/add", `{"a": -0.0, "b": -0.0}`)

	if got := strings.TrimSpace(response.Body.String()); got != `{"result":0}` {
		t.Errorf("body = %s, want {\"result\":0}", got)
	}
}

func TestSubtractAndMultiply(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		want string
	}{
		{name: "subtract", path: "/api/v1/subtract", body: `{"a": 2, "b": 3}`, want: `{"result":-1}`},
		{name: "multiply", path: "/api/v1/multiply", body: `{"a": 2, "b": 3}`, want: `{"result":6}`},
		{name: "subtract negative zero", path: "/api/v1/subtract", body: `{"a": -0.0, "b": 0}`, want: `{"result":0}`},
		{name: "multiply negative zero", path: "/api/v1/multiply", body: `{"a": 0, "b": -1}`, want: `{"result":0}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, tt.path, tt.body)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
			}
			if got := strings.TrimSpace(response.Body.String()); got != tt.want {
				t.Errorf("body = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSubtractAndMultiplyRejectBadInput(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{name: "subtract overflow", path: "/api/v1/subtract", body: `{"a": 1e308, "b": -1e308}`, wantStatus: http.StatusUnprocessableEntity, wantCode: "result_out_of_range"},
		{name: "multiply overflow", path: "/api/v1/multiply", body: `{"a": 1e308, "b": 10}`, wantStatus: http.StatusUnprocessableEntity, wantCode: "result_out_of_range"},
		{name: "subtract malformed JSON", path: "/api/v1/subtract", body: `{"a": 2,`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "multiply malformed JSON", path: "/api/v1/multiply", body: `{"a": 2,`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "subtract missing b", path: "/api/v1/subtract", body: `{"a": 2}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "multiply missing b", path: "/api/v1/multiply", body: `{"a": 2}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "subtract wrong type", path: "/api/v1/subtract", body: `{"a": "2", "b": 3}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "multiply wrong type", path: "/api/v1/multiply", body: `{"a": "2", "b": 3}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "subtract unknown field", path: "/api/v1/subtract", body: `{"a": 2, "b": 3, "c": 4}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "multiply unknown field", path: "/api/v1/multiply", body: `{"a": 2, "b": 3, "c": 4}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, tt.path, tt.body)
			assertError(t, response, tt.wantStatus, tt.wantCode)
		})
	}
}

func TestDivideReturnsResult(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/divide", `{"a": 6, "b": 3}`)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	if got := strings.TrimSpace(response.Body.String()); got != `{"result":2}` {
		t.Errorf("body = %s, want {\"result\":2}", got)
	}
}

func TestDivideByZeroReturnsDivisionByZero(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "zero divisor", body: `{"a": 1, "b": 0}`},
		{name: "zero divided by zero", body: `{"a": 0, "b": 0}`},
		{name: "negative zero divisor", body: `{"a": 1, "b": -0.0}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/divide", tt.body)

			assertError(t, response, http.StatusUnprocessableEntity, "division_by_zero")
			var body errorBody
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("decoding error body %q: %v", response.Body, err)
			}
			if body.Error.Message != "cannot divide by zero" {
				t.Errorf("error message = %q, want %q", body.Error.Message, "cannot divide by zero")
			}
		})
	}
}

func TestDivideOverflowReturnsResultOutOfRange(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/divide", `{"a": 1e308, "b": 1e-308}`)

	assertError(t, response, http.StatusUnprocessableEntity, "result_out_of_range")
}

func TestDivideRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"a": 6,`},
		{name: "missing a", body: `{"b": 3}`},
		{name: "missing b never becomes zero", body: `{"a": 6}`},
		{name: "null b never becomes zero", body: `{"a": 6, "b": null}`},
		{name: "wrong type", body: `{"a": 6, "b": "3"}`},
		{name: "unknown field", body: `{"a": 6, "b": 3, "c": 4}`},
		{name: "number too large for float64", body: `{"a": 1e999, "b": 3}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/divide", tt.body)
			assertError(t, response, http.StatusBadRequest, "invalid_request")
		})
	}
}

func TestUnknownRoutesAndWrongMethodsUseDefaultResponses(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{name: "unknown route", method: http.MethodPost, path: "/api/v1/unknown", wantStatus: http.StatusNotFound},
		{name: "wrong method", method: http.MethodGet, path: "/api/v1/add", wantStatus: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, tt.method, tt.path, "")

			if response.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
				t.Errorf("Content-Type = %q, want text/plain", got)
			}
		})
	}
}

func TestPercentageReturnsResult(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "whole numbers", body: `{"a": 20, "b": 50}`, want: `{"result":10}`},
		{name: "negative percent", body: `{"a": -20, "b": 50}`, want: `{"result":-10}`},
		{name: "decimal percent", body: `{"a": 12.5, "b": 80}`, want: `{"result":10}`},
		{name: "product overflows but the result fits", body: `{"a": 1e308, "b": 50}`, want: `{"result":5e+307}`},
		{name: "negative zero", body: `{"a": 0, "b": -50}`, want: `{"result":0}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/percentage", tt.body)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
			}
			if got := strings.TrimSpace(response.Body.String()); got != tt.want {
				t.Errorf("body = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPercentageOverflowReturnsResultOutOfRange(t *testing.T) {
	response := send(t, http.MethodPost, "/api/v1/percentage", `{"a": 1e308, "b": 1000}`)

	assertError(t, response, http.StatusUnprocessableEntity, "result_out_of_range")
}

func TestPercentageRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"a": 20,`},
		{name: "missing b", body: `{"a": 20}`},
		{name: "string field", body: `{"a": "20", "b": 50}`},
		{name: "unknown field", body: `{"a": 20, "b": 50, "c": 1}`},
		{name: "number too large for float64", body: `{"a": 1e999, "b": 50}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := send(t, http.MethodPost, "/api/v1/percentage", tt.body)
			assertError(t, response, http.StatusBadRequest, "invalid_request")
		})
	}
}

func assertError(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()
	if response.Code != wantStatus {
		t.Errorf("status = %d, want %d; body: %s", response.Code, wantStatus, response.Body)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body errorBody
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding error body %q: %v", response.Body, err)
	}
	if body.Error.Code != wantCode {
		t.Errorf("error code = %q, want %q", body.Error.Code, wantCode)
	}
	if body.Error.Message == "" {
		t.Error("error message is empty")
	}
}
