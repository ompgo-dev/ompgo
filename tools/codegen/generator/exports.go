package generator

import "github.com/ompgo-dev/ompgo/tools/codegen/model"

// GenerateTypes emits the core omp wrapper types.
func (g *Generator) GenerateTypes(outDir string) error {
	return generateTypes(outDir, g.modulePath)
}

// LoadAPIGroups loads api.json into API groups.
func LoadAPIGroups(path string) (map[string][]model.APIFunction, error) {
	return loadAPIGroups(path)
}

// LoadEventGroups loads events.json into event groups.
func LoadEventGroups(path string) (map[string][]model.Event, error) {
	return loadEventGroups(path)
}

// FilterAPIGroups filters API groups by a comma-separated group filter.
func FilterAPIGroups(apiGroups map[string][]model.APIFunction, groupFilter string) map[string][]model.APIFunction {
	return filterAPIGroups(apiGroups, groupFilter)
}

// RemoveExcludedDomains removes generated domain folders not in the allowed set.
func RemoveExcludedDomains(outDir string, allowedGroups map[string][]model.APIFunction) error {
	return removeExcludedDomains(outDir, allowedGroups)
}

// GenerateHelpers emits domain helper wrappers.
func (g *Generator) GenerateHelpers(outDir string, apiGroups map[string][]model.APIFunction) error {
	return generateHelpers(outDir, apiGroups, g.modulePath)
}

// GenerateEvents emits event structs.
func (g *Generator) GenerateEvents(eventGroups map[string][]model.Event, outDir string) error {
	return generateEvents(eventGroups, outDir)
}

// GenerateEventHandlers emits the event handler interface and base type.
func (g *Generator) GenerateEventHandlers(eventGroups map[string][]model.Event, outDir string) error {
	return generateEventHandlers(eventGroups, outDir)
}

// GenerateRuntimeEvents emits runtime event forwarders.
func (g *Generator) GenerateRuntimeEvents(eventGroups map[string][]model.Event, outDir string) error {
	return generateRuntimeEvents(eventGroups, outDir, g.modulePath)
}

// GenerateRuntimeEventExports emits C-exported runtime event handlers.
func (g *Generator) GenerateRuntimeEventExports(eventGroups map[string][]model.Event, outDir string) error {
	return generateRuntimeEventExports(eventGroups, outDir)
}

// GenerateCAPIEventBindings emits the CAPI event registration bindings.
func (g *Generator) GenerateCAPIEventBindings(eventGroups map[string][]model.Event, outDir string) error {
	return generateCAPIEventBindings(eventGroups, outDir)
}

// GenerateCAPIAPIWrappers emits CAPI API wrappers.
func (g *Generator) GenerateCAPIAPIWrappers(apiGroups map[string][]model.APIFunction, outDir string) error {
	return generateCAPIAPIWrappers(apiGroups, outDir, g.modulePath)
}

// GenerateGamemodeBindings emits gamemode API wrappers.
func (g *Generator) GenerateGamemodeBindings(apiGroups map[string][]model.APIFunction, outDir string) error {
	return generateGamemodeBindings(apiGroups, outDir, g.modulePath)
}

// GenerateConstants emits curated omp constants (weapons, keys, etc.).
func (g *Generator) GenerateConstants(outDir string) error {
	return generateConstants(outDir)
}
