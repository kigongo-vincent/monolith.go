package main

import (
	"log"

	"github.com/kigongo-vincent/monolith.go.git/internal/settings"
)

func initSettings() *settings.Config {
	cfg, err := settings.Load(".env")
	if err != nil {
		log.Println("no .env found, using defaults")
		cfg, _ = settings.Load("")
	}
	return cfg
}
