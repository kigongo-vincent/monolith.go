package server

import (
	"testing"

	"github.com/kigongo-vincent/monolith.go.git/pkg/app"
	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
	"github.com/kigongo-vincent/monolith.go.git/pkg/result"
)

func dummyHandler(*app.App, db.DB, *integrations.Integrations) result.Result {
	return result.Ok("ok")
}

func TestRouterFind(t *testing.T) {
	r := &Router{}
	r.Get("/health", dummyHandler)
	r.Get("/foo", dummyHandler)
	if r.find("GET", "/health") == nil {
		t.Fatal("find GET /health should return handler")
	}
	if r.find("GET", "/foo") == nil {
		t.Fatal("find GET /foo should return handler")
	}
	if r.find("POST", "/health") != nil {
		t.Fatal("find POST /health should return nil")
	}
	if r.find("GET", "/missing") != nil {
		t.Fatal("find GET /missing should return nil")
	}
}

func TestRouterAll(t *testing.T) {
	r := &Router{}
	r.All("/api", func() {
		r.Get("/health", dummyHandler)
	})
	if r.find("GET", "/api/health") == nil {
		t.Fatal("find GET /api/health should return handler")
	}
	if r.find("GET", "/health") != nil {
		t.Fatal("find GET /health should return nil after All")
	}
}
