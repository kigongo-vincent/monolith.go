package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const frameworkModule = "github.com/kigongo-vincent/monolith.go.git"

// Run parses args and runs the matching subcommand.
func Run(args []string) error {
	if len(args) < 2 {
		return runHelp()
	}
	cmd := args[1]
	switch cmd {
	case "new":
		return runNew(args[2:])
	case "create":
		return runCreate(args[2:])
	case "sync-client", "benchmark":
		return runStub(cmd)
	default:
		return runHelp()
	}
}

func runHelp() error {
	fmt.Println("Monolith CLI — scaffold Go apps with optional React UI")
	fmt.Println("Usage: monolith <command> [args]")
	fmt.Println("  new <appname>     Scaffold a new app (prompts for React UI)")
	fmt.Println("  create feature   Scaffold a feature (handlers, routes, repo)")
	fmt.Println("  create component Scaffold a React component")
	fmt.Println("  sync-client      Sync backend types to app/src/types.ts")
	fmt.Println("  benchmark        Run benchmark suite")
	return nil
}

func runNew(args []string) error {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	withUI := fs.Bool("ui", false, "add React UI")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: monolith new [--ui] <appname>")
	}
	appName := strings.TrimSpace(fs.Arg(0))
	if appName == "" {
		return fmt.Errorf("appname cannot be empty")
	}
	if !*withUI {
		var err error
		*withUI, err = promptAddReactUI()
		if err != nil {
			return err
		}
	}
	return Scaffold(appName, *withUI, appName)
}

func promptAddReactUI() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Add React UI? (y/n): ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}

func runCreate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: monolith create <feature|component> [name]")
	}
	switch args[0] {
	case "feature", "component":
		return runStub("create " + args[0])
	default:
		return runHelp()
	}
}

func runStub(name string) error {
	fmt.Printf("%s: not implemented yet\n", name)
	return nil
}
