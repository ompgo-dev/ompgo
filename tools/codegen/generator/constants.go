package generator

import "strings"

const (
	DefaultModulePath = "github.com/ompgo-dev/ompgo"
)

// Generator holds immutable configuration for a code generation run.
type Generator struct {
	modulePath string
}

// New returns a generator configured for a single code generation run.
func New(modulePath string) *Generator {
	return &Generator{modulePath: normalizeModulePath(modulePath)}
}

func normalizeModulePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultModulePath
	}
	return value
}
