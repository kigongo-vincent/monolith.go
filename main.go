package main

import (
	"log"

	_ "modernc.org/sqlite"

	"github.com/kigongo-vincent/monolith.go.git/internal/app"
	"github.com/kigongo-vincent/monolith.go.git/internal/db"
	"github.com/kigongo-vincent/monolith.go.git/internal/integrations"
	"github.com/kigongo-vincent/monolith.go.git/internal/result"
	"github.com/kigongo-vincent/monolith.go.git/internal/server"
)

func main() {
	cfg := initSettings()
	database, err := db.New(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()
	storage, _ := integrations.NewLocalStorage(cfg.Storage.LocalPath)
	ints := integrations.New(storage, nil, nil, nil)
	registerRoutes()
	server.Run(cfg.Server.Port, database, ints)
}

func registerRoutes() {
	server.Get("/health", healthHandler)
	server.All("/api", func() {
		server.Get("/health", healthHandler)
	})
}

func healthHandler(a *app.App, database db.DB, ints *integrations.Integrations) result.Result {
	return result.Ok("ok")
}
