package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/gamemode.go.tmpl
var templates embed.FS

type gamemodeTemplateData struct {
	Name     string
	TypeName string
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		cmdInit(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "ompgo init -name <gamemode> -module <module/path> [-out <dir>]")
}

func cmdInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	name := fs.String("name", "", "Gamemode name (used for folder and component name)")
	module := fs.String("module", "", "Go module path (required)")
	outDir := fs.String("out", "", "Output directory (default: ./<name>)")
	_ = fs.Parse(args)

	if strings.TrimSpace(*name) == "" {
		fatal(errors.New("-name is required"))
	}
	if strings.TrimSpace(*module) == "" {
		fatal(errors.New("-module is required"))
	}

	cleanName := normalizeName(*name)
	if cleanName == "" {
		fatal(errors.New("invalid -name"))
	}
	cleanModule := strings.TrimSpace(*module)
	if cleanModule == "" {
		fatal(errors.New("invalid -module"))
	}

	typeName := pascalCase(cleanName)
	if typeName == "" {
		fatal(errors.New("invalid -name for type"))
	}

	outPath := *outDir
	if outPath == "" {
		outPath = cleanName
	}
	outPath = filepath.Clean(outPath)
	if _, err := os.Stat(outPath); err == nil {
		fatal(fmt.Errorf("output already exists: %s", outPath))
	}
	if err := os.MkdirAll(outPath, 0o755); err != nil {
		fatal(err)
	}

	if err := writeGoMod(outPath, cleanModule); err != nil {
		fatal(err)
	}
	if err := writeMain(outPath, gamemodeTemplateData{Name: cleanName, TypeName: typeName}); err != nil {
		fatal(err)
	}

	fmt.Printf("Created gamemode at %s\n", filepath.Join(outPath, "main.go"))
}

func writeGoMod(outPath, module string) error {
	content := fmt.Sprintf("module %s\n\ngo 1.21\n", module)
	return os.WriteFile(filepath.Join(outPath, "go.mod"), []byte(content), 0o644)
}

func writeMain(outPath string, data gamemodeTemplateData) error {
	b, err := templates.ReadFile("templates/gamemode.go.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("gamemode").Parse(string(b))
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(outPath, "main.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	return strings.Trim(name, "-_")
}

func pascalCase(s string) string {
	var b strings.Builder
	upperNext := true
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			upperNext = true
			continue
		}
		if upperNext {
			b.WriteRune(unicode.ToUpper(r))
			upperNext = false
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
