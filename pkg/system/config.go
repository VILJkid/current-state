package system

import (
	"fmt"
	"path/filepath"
)

// AllowedPaths defines filesystem paths that are safe to query
var AllowedPaths = map[string]bool{
	"/":     true, // Root filesystem
	"/home": true, // User home directories
	"/tmp":  true, // Temporary files
}

// ValidatePath ensures the requested path is in the whitelist.
// Returns an error if the path is not allowed.
func ValidatePath(path string) error {
	// Normalize the path to prevent bypass attempts
	normalized := filepath.Clean(path)

	if !AllowedPaths[normalized] {
		return fmt.Errorf("access to path %q not allowed", path)
	}
	return nil
}

// AddAllowedPath adds a new path to the whitelist.
// Should only be called during initialization.
func AddAllowedPath(path string) {
	normalized := filepath.Clean(path)
	AllowedPaths[normalized] = true
}
