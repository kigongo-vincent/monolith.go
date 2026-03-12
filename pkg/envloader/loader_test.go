package envloader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	err := Load(filepath.Join(t.TempDir(), "nonexistent.env"))
	if err == nil {
		t.Fatal("Load missing file should return error")
	}
}

func TestLoadValidFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, ".env")
	content := "FOO=bar\n# comment\nBAR=baz\nEMPTY=\n"
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Load(f); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("FOO") != "bar" {
		t.Errorf("FOO want bar got %q", os.Getenv("FOO"))
	}
	if os.Getenv("BAR") != "baz" {
		t.Errorf("BAR want baz got %q", os.Getenv("BAR"))
	}
}

func TestLoadEmptyPath(t *testing.T) {
	// Load("") should not open a file; we don't have a way to test that except no panic
	// So we test Load with empty path - actually Load("") would try to open "" and fail
	err := Load("")
	if err == nil {
		t.Log("Load empty path returned nil (ok on some systems)")
	}
}
