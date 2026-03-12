package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kigongo-vincent/monolith.go.git/pkg/app"
	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
	"github.com/kigongo-vincent/monolith.go.git/pkg/result"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", "/"},
		{"/", "/"},
		{"/health", "/health"},
		{"/health/", "/health"},
		{"health", "/health"},
	}
	for _, tt := range tests {
		got := normalizePath(tt.in)
		if got != tt.want {
			t.Errorf("normalizePath(%q) want %q got %q", tt.in, tt.want, got)
		}
	}
}

func TestServeHTTP200(t *testing.T) {
	router := &Router{}
	router.Get("/health", func(*app.App, db.DB, *integrations.Integrations) result.Result {
		return result.Ok("ok")
	})
	srv := New(router, nil, integrations.New(nil, nil, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("status want 200 got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("body want ok got %q", rec.Body.String())
	}
}

func TestServeHTTP404(t *testing.T) {
	router := &Router{}
	srv := New(router, nil, integrations.New(nil, nil, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != 404 {
		t.Errorf("status want 404 got %d", rec.Code)
	}
}

func TestServeHTTPErr(t *testing.T) {
	router := &Router{}
	router.Get("/err", func(*app.App, db.DB, *integrations.Integrations) result.Result {
		return result.Err(401, "Unauthorized")
	})
	srv := New(router, nil, integrations.New(nil, nil, nil, nil))
	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != 401 {
		t.Errorf("status want 401 got %d", rec.Code)
	}
}
