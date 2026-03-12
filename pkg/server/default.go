package server

import (
	"log"

	"github.com/kigongo-vincent/monolith.go.git/pkg/app"
	"github.com/kigongo-vincent/monolith.go.git/pkg/db"
	"github.com/kigongo-vincent/monolith.go.git/pkg/integrations"
	"github.com/kigongo-vincent/monolith.go.git/pkg/result"
	"github.com/kigongo-vincent/monolith.go.git/pkg/settings"
)

// RunDefault loads config, opens DB, wires integrations, registers default routes, and runs the server.
func RunDefault() {
	cfg := loadConfig()
	database, err := db.New(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()
	storage, _ := integrations.NewLocalStorage(cfg.Storage.LocalPath)
	ints := integrations.New(storage, nil, nil, nil)
	registerDefaultRoutes()
	Run(cfg.Server.Port, database, ints)
}

func loadConfig() *settings.Config {
	cfg, err := settings.Load(".env")
	if err != nil {
		log.Println("no .env found, using defaults")
		cfg, _ = settings.Load("")
	}
	return cfg
}

func registerDefaultRoutes() {
	Get("/health", healthHandler)
	All("/api", func() {
		Get("/health", healthHandler)
	})
}

func healthHandler(a *app.App, _ db.DB, _ *integrations.Integrations) result.Result {
	return result.Ok("ok")
}
