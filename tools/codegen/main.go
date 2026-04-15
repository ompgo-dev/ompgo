package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ompgo-dev/ompgo/tools/codegen/generator"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	apiPath := flag.String("api", "tools/codegen/data/openmp-capi/api.json", "Path to api.json")
	eventsPath := flag.String("events", "tools/codegen/data/openmp-capi/events.json", "Path to events.json")
	outDir := flag.String("out", "pkg/omp", "Output directory")
	groupFilter := flag.String("groups", "", "Comma-separated list of API groups to generate (default: minimal gamemode set)")
	modulePath := flag.String("repo", generator.DefaultModulePath, "Go module import path for generated code")
	flag.Parse()

	gen := generator.New(*modulePath)

	fmt.Println("OMP Go Code Generator")
	fmt.Println("=====================")

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", *outDir, err)
	}

	if err := gen.GenerateTypes(*outDir); err != nil {
		return fmt.Errorf("generate types: %w", err)
	}

	apiGroups, err := generator.LoadAPIGroups(*apiPath)
	if err != nil {
		return fmt.Errorf("load API groups from %s: %w", *apiPath, err)
	}

	apiGroups = generator.FilterAPIGroups(apiGroups, *groupFilter)
	if err := generator.RemoveExcludedDomains(*outDir, apiGroups); err != nil {
		return fmt.Errorf("remove excluded domains: %w", err)
	}

	if err := gen.GenerateHelpers(*outDir, apiGroups); err != nil {
		return fmt.Errorf("generate helpers: %w", err)
	}

	eventGroups, err := generator.LoadEventGroups(*eventsPath)
	if err != nil {
		return fmt.Errorf("load event groups from %s: %w", *eventsPath, err)
	}

	if err := gen.GenerateEvents(eventGroups, *outDir); err != nil {
		return fmt.Errorf("generate events: %w", err)
	}

	if err := gen.GenerateEventHandlers(eventGroups, *outDir); err != nil {
		return fmt.Errorf("generate event handlers: %w", err)
	}

	// Resolve output roots from outDir
	pkgDir := filepath.Dir(*outDir)
	runtimeDir := filepath.Join(pkgDir, "runtime")

	// Generate runtime event forwarders
	if err := gen.GenerateRuntimeEvents(eventGroups, runtimeDir); err != nil {
		return fmt.Errorf("generate runtime events: %w", err)
	}

	// Generate C exports for events (runtime package)
	if err := gen.GenerateRuntimeEventExports(eventGroups, runtimeDir); err != nil {
		return fmt.Errorf("generate runtime exports: %w", err)
	}

	// Generate CAPI event registration bindings
	if err := gen.GenerateCAPIEventBindings(eventGroups, runtimeDir); err != nil {
		return fmt.Errorf("generate CAPI event bindings: %w", err)
	}

	// Generate CAPI API wrappers (C + Go)
	if err := gen.GenerateCAPIAPIWrappers(apiGroups, runtimeDir); err != nil {
		return fmt.Errorf("generate CAPI API wrappers: %w", err)
	}

	// Generate curated omp constants (weapons, keys, etc.)
	ompDir := filepath.Join(pkgDir, "omp")
	if err := gen.GenerateConstants(ompDir); err != nil {
		return fmt.Errorf("generate constants: %w", err)
	}

	gamemodeDir := filepath.Join(pkgDir, "gamemode")
	// Generate gamemode API wrappers
	if err := gen.GenerateGamemodeBindings(apiGroups, gamemodeDir); err != nil {
		return fmt.Errorf("generate gamemode API wrappers: %w", err)
	}

	fmt.Println("Done!")
	return nil
}
