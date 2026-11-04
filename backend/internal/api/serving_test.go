package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fullstack-calculator/backend/internal/api"
)

func request(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}

// staticDir creates a built-frontend directory. Besides the real files, it holds
// files at paths the API owns, which must never be served.
func staticDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"index.html":     "<h1>calculator</h1>",
		"assets/app.js":  "console.log('app')",
		"healthz":        "static healthz",
		"api/v1/unknown": "static unknown",
		"api/v1/add":     "static add",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestServesStaticFiles(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantBody        string
		wantContentType string
	}{
		{name: "root serves index.html", path: "/", wantBody: "<h1>calculator</h1>", wantContentType: "text/html"},
		{name: "nested asset", path: "/assets/app.js", wantBody: "console.log('app')", wantContentType: "text/javascript"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := request(api.NewHandler(staticDir(t)), http.MethodGet, tt.path)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
			}
			if got := response.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, tt.wantContentType) {
				t.Errorf("Content-Type = %q, want %s", got, tt.wantContentType)
			}
		})
	}
}

func TestHealthCheckReturnsOK(t *testing.T) {
	for _, mode := range modes(t) {
		t.Run(mode.name, func(t *testing.T) {
			response := request(mode.handler, http.MethodGet, "/healthz")

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
			}
			if got := response.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if got := strings.TrimSpace(response.Body.String()); got != `{"status":"ok"}` {
				t.Errorf("body = %s, want {\"status\":\"ok\"}", got)
			}
		})
	}
}

func TestCalculatorRoutesWorkInBothModes(t *testing.T) {
	for _, mode := range modes(t) {
		t.Run(mode.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			body := strings.NewReader(`{"a": 2, "b": 3}`)
			mode.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/add", body))

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", recorder.Code, http.StatusOK, recorder.Body)
			}
			if got := strings.TrimSpace(recorder.Body.String()); got != `{"result":5}` {
				t.Errorf("body = %s, want {\"result\":5}", got)
			}
		})
	}
}

func TestStaticServingDoesNotChangeDefaultResponses(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantAllow  string
	}{
		{name: "GET on a calculator route", method: http.MethodGet, path: "/api/v1/add", wantStatus: http.StatusMethodNotAllowed, wantAllow: "POST"},
		{name: "PUT on a calculator route", method: http.MethodPut, path: "/api/v1/add", wantStatus: http.StatusMethodNotAllowed, wantAllow: "POST"},
		{name: "POST on the health check", method: http.MethodPost, path: "/healthz", wantStatus: http.StatusMethodNotAllowed, wantAllow: "GET, HEAD"},
		{name: "POST to an unknown API route", method: http.MethodPost, path: "/api/v1/unknown", wantStatus: http.StatusNotFound},
		{name: "GET on an unknown API route", method: http.MethodGet, path: "/api/v1/unknown", wantStatus: http.StatusNotFound},
		{name: "GET on the API root", method: http.MethodGet, path: "/api", wantStatus: http.StatusNotFound},
		{name: "POST to the API root", method: http.MethodPost, path: "/api", wantStatus: http.StatusNotFound},
		{name: "GET on the API root with a slash", method: http.MethodGet, path: "/api/", wantStatus: http.StatusNotFound},
		{name: "GET on an unknown path", method: http.MethodGet, path: "/missing", wantStatus: http.StatusNotFound},
		{name: "POST to an unknown path", method: http.MethodPost, path: "/missing", wantStatus: http.StatusNotFound},
		{name: "PUT on an unknown nested path", method: http.MethodPut, path: "/assets/missing.js", wantStatus: http.StatusNotFound},
	}

	for _, mode := range modes(t) {
		for _, tt := range tests {
			t.Run(mode.name+"/"+tt.name, func(t *testing.T) {
				response := request(mode.handler, tt.method, tt.path)

				if response.Code != tt.wantStatus {
					t.Errorf("status = %d, want %d; body: %s", response.Code, tt.wantStatus, response.Body)
				}
				if got := response.Header().Get("Allow"); got != tt.wantAllow {
					t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
				}
				if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
					t.Errorf("Content-Type = %q, want text/plain", got)
				}
			})
		}
	}
}

func TestAPIOnlyModeServesNoFrontend(t *testing.T) {
	response := request(api.NewHandler(""), http.MethodGet, "/")

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d; body: %s", response.Code, http.StatusNotFound, response.Body)
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
		t.Errorf("Content-Type = %q, want text/plain", got)
	}
}

func TestStaticFilesRejectOtherMethods(t *testing.T) {
	for _, path := range []string{"/", "/assets/app.js", "/assets/"} {
		t.Run(path, func(t *testing.T) {
			response := request(api.NewHandler(staticDir(t)), http.MethodPost, path)

			if response.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want %d; body: %s", response.Code, http.StatusMethodNotAllowed, response.Body)
			}
			if got := response.Header().Get("Allow"); got != "GET, HEAD" {
				t.Errorf("Allow = %q, want GET, HEAD", got)
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
				t.Errorf("Content-Type = %q, want text/plain", got)
			}
		})
	}
}

type mode struct {
	name    string
	handler http.Handler
}

func modes(t *testing.T) []mode {
	t.Helper()
	return []mode{
		{name: "API-only", handler: api.NewHandler("")},
		{name: "static", handler: api.NewHandler(staticDir(t))},
	}
}
