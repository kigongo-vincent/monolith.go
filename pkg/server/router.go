package server

import (
	"github.com/kigongo-vincent/monolith.go.git/pkg/app"
	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
	"github.com/kigongo-vincent/monolith.go.git/pkg/result"
)

// Handler is the standard handler signature.
type Handler func(*app.App, db.DB, *integrations.Integrations) result.Result

type route struct {
	method  string
	path    string
	handler Handler
}

// Router holds routes and prefix for grouping.
type Router struct {
	routes []route
	prefix string
}

// Get registers a GET handler.
func (r *Router) Get(path string, h Handler) {
	r.routes = append(r.routes, route{"GET", r.prefix + path, h})
}

// Post registers a POST handler.
func (r *Router) Post(path string, h Handler) {
	r.routes = append(r.routes, route{"POST", r.prefix + path, h})
}

// All sets a path prefix and runs fn to register nested routes.
func (r *Router) All(prefix string, fn func()) {
	old := r.prefix
	r.prefix = old + prefix
	fn()
	r.prefix = old
}

func (r *Router) find(method, path string) Handler {
	for i := len(r.routes) - 1; i >= 0; i-- {
		ro := r.routes[i]
		if ro.method == method && ro.path == path {
			return ro.handler
		}
	}
	return nil
}
