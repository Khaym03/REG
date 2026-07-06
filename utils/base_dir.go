package utils

import (
	"os"
	"path/filepath"
)

// Prefer per-user writable directory instead of current working directory
func BaseDir() string {
	if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
		return filepath.Join(dir, "REG")
	}

	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "REG")
	}

	if dir, err := os.UserHomeDir(); err == nil && dir != "" {
		return filepath.Join(dir, ".local", "share", "REG")
	}

	// Fallback to working directory if all else fails
	wd, err := os.Getwd()
	if err != nil {
		return "REG"
	}

	return filepath.Join(wd, "REG")
}
