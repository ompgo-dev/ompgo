package generator

import (
	"fmt"
	"go/format"
	"os"
)

func writeFormattedGoFile(path string, src []byte) error {
	formatted, err := format.Source(src)
	if err != nil {
		// If formatting fails, still write the file to help debugging.
		_ = os.WriteFile(path, src, 0o644)
		return fmt.Errorf("format %s: %w", path, err)
	}
	return os.WriteFile(path, formatted, 0o644)
}

