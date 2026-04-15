package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

func TestGenerateRuntimeEventsUsesConfiguredRepoPath(t *testing.T) {
	gen := New("example.com/custom/ompgo")
	outDir := t.TempDir()
	eventGroups := map[string][]model.Event{
		"Player": {
			{Name: "onPlayerConnect", BadRet: "none"},
		},
	}

	if err := gen.GenerateRuntimeEvents(eventGroups, outDir); err != nil {
		t.Fatalf("GenerateRuntimeEvents() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outDir, "events_gen.go"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if !strings.Contains(string(content), "\"example.com/custom/ompgo/pkg/omp\"") {
		t.Fatalf("generated runtime events did not use configured repo path:\n%s", content)
	}
	if !strings.Contains(string(content), "loadHandlerSnapshot(&onPlayerConnectHandlers)") {
		t.Fatalf("generated runtime events did not use handler snapshots:\n%s", content)
	}
	if !strings.Contains(string(content), "cfg := currentEventDispatchConfig()") {
		t.Fatalf("generated runtime events did not preload dispatch config:\n%s", content)
	}
	if strings.Contains(string(content), "defer cancel()") {
		t.Fatalf("generated runtime events still defer an event cancel function:\n%s", content)
	}
	if strings.Contains(string(content), "eventHandlersMu.RLock()") {
		t.Fatalf("generated runtime events still hold a read lock while dispatching:\n%s", content)
	}
}
