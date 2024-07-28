package utils

import (
	"os"
	"path/filepath"
)

// Shared utility functions
func FindProjectRoot(marker string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		// Check if the marker exists in this directory
		path := filepath.Join(dir, marker)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return dir
		}

		// Move to the parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// We've reached the root directory and didn't find the marker
			return ""
		}
		dir = parentDir
	}
}

func GetEnviroment() string {
	if env := os.Getenv("ENV"); env != "" {
		return env
	} else {
		return "local"
	}
}
