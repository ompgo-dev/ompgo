package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

func TestEventGeneratorsUseSortedGroupOrder(t *testing.T) {
	eventGroups := map[string][]model.Event{
		"Vehicle": {
			{Name: "onVehicleSpawn", BadRet: "none"},
		},
		"Actor": {
			{Name: "onActorStreamIn", BadRet: "none"},
		},
	}

	gen := New("example.com/custom/ompgo")
	ompOutDir := t.TempDir()
	if err := generateEvents(eventGroups, ompOutDir); err != nil {
		t.Fatalf("generateEvents() error = %v", err)
	}
	assertStringOrderInFile(t, filepath.Join(ompOutDir, "events_gen.go"), "ActorStreamInEvent", "VehicleSpawnEvent")

	if err := generateEventHandlers(eventGroups, ompOutDir); err != nil {
		t.Fatalf("generateEventHandlers() error = %v", err)
	}
	assertStringOrderInFile(t, filepath.Join(ompOutDir, "event_handlers_gen.go"), "OnActorStreamIn(ctx context.Context", "OnVehicleSpawn(ctx context.Context")

	runtimeOutDir := t.TempDir()
	if err := gen.GenerateRuntimeEvents(eventGroups, runtimeOutDir); err != nil {
		t.Fatalf("GenerateRuntimeEvents() error = %v", err)
	}
	assertStringOrderInFile(t, filepath.Join(runtimeOutDir, "events_gen.go"), "type onActorStreamInHandlersFunc", "type onVehicleSpawnHandlersFunc")

	if err := gen.GenerateRuntimeEventExports(eventGroups, runtimeOutDir); err != nil {
		t.Fatalf("GenerateRuntimeEventExports() error = %v", err)
	}
	assertStringOrderInFile(t, filepath.Join(runtimeOutDir, "events_exports_gen.go"), "func OMPGO_OnActorStreamIn", "func OMPGO_OnVehicleSpawn")

	if err := gen.GenerateCAPIEventBindings(eventGroups, runtimeOutDir); err != nil {
		t.Fatalf("GenerateCAPIEventBindings() error = %v", err)
	}
	assertStringOrderInFile(t, filepath.Join(runtimeOutDir, "capi_events_cgo_gen.go"), "extern void OMPGO_OnActorStreamIn", "extern void OMPGO_OnVehicleSpawn")
}

func TestGenerateTypesUsesStableColorOrder(t *testing.T) {
	gen := New("example.com/custom/ompgo")
	outDir := t.TempDir()
	if err := gen.GenerateTypes(outDir); err != nil {
		t.Fatalf("GenerateTypes() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outDir, "types_gen.go"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	orderedColors := []string{
		"ColorWhite",
		"ColorBlack",
		"ColorRed",
		"ColorGreen",
		"ColorBlue",
		"ColorYellow",
		"ColorOrange",
		"ColorPurple",
		"ColorGrey",
	}
	assertStringsInOrder(t, string(content), orderedColors...)
}

func assertStringOrderInFile(t *testing.T, path string, ordered ...string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	assertStringsInOrder(t, string(content), ordered...)
}

func assertStringsInOrder(t *testing.T, content string, ordered ...string) {
	t.Helper()
	last := -1
	for _, needle := range ordered {
		index := strings.Index(content, needle)
		if index == -1 {
			t.Fatalf("expected to find %q in generated content", needle)
		}
		if index <= last {
			t.Fatalf("expected %q to appear after the previous marker in generated content", needle)
		}
		last = index
	}
}
