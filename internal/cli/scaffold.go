package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Scaffold creates a new app in destDir with optional React UI.
func Scaffold(appName string, withUI bool, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", destDir, err)
	}
	cwd, _ := os.Getwd()
	appAbs := filepath.Join(cwd, destDir)
	replacePath, err := replacePathToFramework(appAbs)
	if err != nil {
		return err
	}
	if err := writeGoFiles(destDir, appName, replacePath); err != nil {
		return err
	}
	if withUI {
		if err := writeAppFiles(destDir); err != nil {
			return err
		}
	}
	fmt.Printf("Scaffolded app %q in %s\n", appName, destDir)
	fmt.Println("Next: cd", destDir, "&& go mod tidy && go run .")
	return nil
}

// replacePathToFramework returns a relative path from the app dir to the framework root (where go.mod has the framework module).
func replacePathToFramework(appDir string) (string, error) {
	dir := filepath.Dir(appDir)
	for {
		gomod := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(gomod)
		if err == nil && strings.Contains(string(data), frameworkModule) {
			rel, err := filepath.Rel(appDir, dir)
			if err != nil {
				return "", err
			}
			if rel == "" {
				rel = "."
			}
			return rel, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("framework root not found: run monolith new from the framework repo or a subdir of it")
		}
		dir = parent
	}
}

func writeGoFiles(destDir, module, replacePath string) error {
	if err := writeFile(destDir, "main.go", buildMainGoContent()); err != nil {
		return err
	}
	if err := writeGoMod(destDir, module, replacePath); err != nil {
		return err
	}
	return writeFile(destDir, ".env.example", envExampleContent)
}

func buildMainGoContent() string {
	return `package main

import (
	_ "modernc.org/sqlite"

	"` + frameworkModule + `/pkg/server"
)

func main() {
	server.RunDefault()
}
`
}

func writeGoMod(destDir, module, replacePath string) error {
	t, err := template.New("go.mod").Parse(goModContent)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(destDir, "go.mod"))
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, map[string]string{"Module": module, "ReplacePath": replacePath})
}

func writeFile(dir, name, content string) error {
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
}

const goModContent = `module {{.Module}}

go 1.24.1

replace ` + frameworkModule + ` => {{.ReplacePath}}

require ` + frameworkModule + ` v0.0.0
`

const envExampleContent = `PORT=8080
DB_DRIVER=sqlite
DB_DSN=file:monolith.db
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=storage
PAYMENT_PROVIDER=stripe
CACHE_BACKEND=memory
JWT_SECRET=change-me-in-production
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
`
