package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldCreatesFiles(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold("testapp", false, dir)
	if err != nil {
		t.Fatal(err)
	}
	wantFiles := []string{"main.go", "settings.go", "go.mod", ".env.example"}
	for _, name := range wantFiles {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("missing file %s", name)
		}
	}
}

func TestScaffoldMainGoContent(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold("testapp", false, dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(b)
	if len(content) == 0 {
		t.Fatal("main.go should not be empty")
	}
	if !strings.Contains(content, "pkg/app") || !strings.Contains(content, "pkg/server") {
		t.Error("main.go should import pkg/app and pkg/server")
	}
}

func TestScaffoldWithUI(t *testing.T) {
	dir := t.TempDir()
	err := Scaffold("testapp", true, dir)
	if err != nil {
		t.Fatal(err)
	}
	appDir := filepath.Join(dir, "app")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		t.Fatal("app/ directory should exist when withUI=true")
	}
	packageJSON := filepath.Join(appDir, "package.json")
	if _, err := os.Stat(packageJSON); os.IsNotExist(err) {
		t.Error("app/package.json should exist")
	}
}
