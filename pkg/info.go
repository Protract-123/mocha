package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetActiveVersion(packageName string, mochaDir string) (string, error) {
	currentDir := filepath.Join(mochaDir, "apps", packageName, "current")

	currentTarget, err := os.Readlink(currentDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to read junction target of %s: %w", currentDir, err)
	}

	return filepath.Base(currentTarget), nil
}
