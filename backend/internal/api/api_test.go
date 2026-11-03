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
