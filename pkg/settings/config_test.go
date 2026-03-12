package settings

import (
	"os"
	"testing"
)

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("default port want 8080 got %s", cfg.Server.Port)
	}
	if cfg.DB.Driver != "sqlite" {
		t.Errorf("default driver want sqlite got %s", cfg.DB.Driver)
	}
}

func TestGetAfterLoad(t *testing.T) {
	cfg, _ := Load("")
	got := Get()
	if got != cfg {
		t.Fatal("Get() should return same config as Load")
	}
}

func TestLoadRespectsEnv(t *testing.T) {
	prev := os.Getenv("PORT")
	os.Setenv("PORT", "3000")
	defer func() {
		if prev == "" {
			os.Unsetenv("PORT")
		} else {
			os.Setenv("PORT", prev)
		}
	}()
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "3000" {
		t.Errorf("port from env want 3000 got %s", cfg.Server.Port)
	}
}
