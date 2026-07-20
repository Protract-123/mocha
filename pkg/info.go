package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/manifest"
)

type Package struct {
	manifest.Info
	Json map[string]any
	Arch string
}

type InstallInfo struct {
	Bucket  string `json:"bucket"`
	Version string `json:"version"`
	Arch    string `json:"architecture"`

	ManifestJson map[string]any `json:"-"`
}

func GetInstallInfo(appDir string) (*InstallInfo, error) {
	installJson, err := os.ReadFile(filepath.Join(appDir, "install.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read install JSON: %w", err)
	}

	installInfo := InstallInfo{}
	if err := json.Unmarshal(installJson, &installInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal install JSON: %w", err)
	}

	manifestJson, err := manifest.GetJson(filepath.Join(appDir, "manifest.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest JSON: %w", err)
	}

	installInfo.ManifestJson = manifestJson
	return &installInfo, nil
}

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
