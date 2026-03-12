package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kigongo-vincent/monolith.go.git/pkg/app"
	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
	"github.com/kigongo-vincent/monolith.go.git/pkg/result"
)

// Server holds router, db, and integrations.
type Server struct {
	router *Router
	db     db.DB
	ints   *integrations.Integrations
}

// New creates a Server with the given router, db, and integrations.
func New(router *Router, database db.DB, ints *integrations.Integrations) *Server {
	return &Server{router: router, db: database, ints: ints}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := normalizePath(r.URL.Path)
	handler := s.router.find(r.Method, path)
	if handler == nil {
		writeResult(w, result.Err(404, "Not Found"))
		return
	}
	a := app.NewApp(r, w)
	res := handler(a, s.db, s.ints)
	writeResult(w, res)
}

func normalizePath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return "/"
	}
	return p
}

func writeResult(w http.ResponseWriter, res result.Result) {
	w.WriteHeader(res.Status())
	if res.IsError() {
		writeJSON(w, map[string]string{"error": res.Body().(string)})
		return
	}
	body := res.Body()
	if body == nil {
		return
	}
	switch v := body.(type) {
	case string:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(v))
	default:
		writeJSON(w, body)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
