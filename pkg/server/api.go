package server

import (
	"log"
	"net/http"

	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
)

var defaultRouter = &Router{}

// Get registers a GET handler on the default router.
func Get(path string, h Handler) {
	defaultRouter.Get(path, h)
}

// Post registers a POST handler on the default router.
func Post(path string, h Handler) {
	defaultRouter.Post(path, h)
}

// All registers a path prefix and runs fn to add nested routes.
func All(prefix string, fn func()) {
	defaultRouter.All(prefix, fn)
}

// Run starts the HTTP server using the default router, db, and integrations.
func Run(port string, database db.DB, ints *integrations.Integrations) {
	srv := New(defaultRouter, database, ints)
	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server: %v", err)
	}
}
