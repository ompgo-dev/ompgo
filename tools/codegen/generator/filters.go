package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

var defaultAllowedGroups = []string{
	"All",
	"Core",
	"Player",
	"Vehicle",
	"Class",
	"Dialog",
	"Checkpoint",
	"RaceCheckpoint",
	"Object",
	"PlayerObject",
	"TextDraw",
	"PlayerTextDraw",
	"TextLabel",
	"PlayerTextLabel",
	"Menu",
	"Pickup",
	"GangZone",
	"Actor",
	"Config",
}

func filterAPIGroups(apiGroups map[string][]model.APIFunction, groupFilter string) map[string][]model.APIFunction {
	allowed := buildAllowedGroups(apiGroups, groupFilter)
	filtered := make(map[string][]model.APIFunction)
	for name, group := range apiGroups {
		if allowed[name] {
			filtered[name] = group
		}
	}
	return filtered
}

func buildAllowedGroups(apiGroups map[string][]model.APIFunction, groupFilter string) map[string]bool {
	allowed := make(map[string]bool)
	if strings.TrimSpace(groupFilter) == "" {
		for _, name := range defaultAllowedGroups {
			if real := matchGroupName(apiGroups, name); real != "" {
				allowed[real] = true
			}
		}
		return allowed
	}

	parts := strings.FieldsFunc(groupFilter, func(r rune) bool {
		return r == ',' || r == ';' || r == '|' || r == ' '
	})
	for _, name := range parts {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if real := matchGroupName(apiGroups, name); real != "" {
			allowed[real] = true
		} else {
			fmt.Printf("Warning: unknown API group '%s' in -groups\n", name)
		}
	}

	return allowed
}

func matchGroupName(apiGroups map[string][]model.APIFunction, name string) string {
	for group := range apiGroups {
		if strings.EqualFold(group, name) {
			return group
		}
	}
	return ""
}

func removeExcludedDomains(outDir string, allowedGroups map[string][]model.APIFunction) error {
	allowedDirs := map[string]bool{}
	for group := range allowedGroups {
		allowedDirs[domainPackageName(group)] = true
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if allowedDirs[name] {
			continue
		}
		if err := os.RemoveAll(filepath.Join(outDir, name)); err != nil {
			return err
		}
	}
	return nil
}
