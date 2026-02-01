package system

import (
	"errors"
	"fmt"
	"syscall"
)

// SanitizeError maps OS errors to user-friendly messages.
// This prevents information leakage about system configuration.
func SanitizeError(err error) error {
	if err == nil {
		return nil
	}

	// Map specific OS errors to generic messages
	if errors.Is(err, syscall.EACCES) {
		return fmt.Errorf("permission denied: insufficient access to system data")
	}
	if errors.Is(err, syscall.ENOENT) {
		return fmt.Errorf("filesystem not found: path may not exist on this system")
	}
	if errors.Is(err, syscall.EIO) {
		return fmt.Errorf("I/O error: unable to read filesystem data")
	}

	// Generic fallback for unmapped errors
	return fmt.Errorf("unable to retrieve system information")
}
