package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/shim"
)

func UnlinkApp(appName string, mochaDir string) error {
	currentDir := filepath.Join(mochaDir, "apps", appName, "current")

	manifestJson, err := manifest.GetJson(filepath.Join(currentDir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("failed to get manifest JSON: %w", err)
	}

	installJson, err := os.ReadFile(filepath.Join(currentDir, "install.json"))
	if err != nil {
		return fmt.Errorf("failed to read install JSON: %w", err)
	}

	installInfo := InstallInfo{}
	if err := json.Unmarshal(installJson, &installInfo); err != nil {
		return fmt.Errorf("failed to unmarshal install JSON: %w", err)
	}

	binaries, err := manifest.GetExecutableEntries(manifestJson, installInfo.Arch)
	if err != nil {
		return fmt.Errorf("failed to get shims to remove: %w", err)
	}

	for _, binary := range binaries {
		shimName := strings.TrimSuffix(filepath.Base(binary.Alias), filepath.Ext(binary.Alias))
		if err := shim.DeleteShim(shimName, mochaDir); err != nil {
			return fmt.Errorf("failed to remove shim %q: %w", shimName, err)
		}
	}

	shortcutDirectory := os.ExpandEnv(filepath.Join("$APPDATA", "Microsoft", "Windows", "Start Menu", "Programs", "Mocha Apps"))

	shortcutEntries := manifest.GetShortcutEntries(manifestJson, installInfo.Arch)
	for _, shortcutEntry := range shortcutEntries {
		shortcutPath := filepath.Join(shortcutDirectory, fmt.Sprintf("%s.lnk", shortcutEntry.Name))

		if err := os.Remove(shortcutPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to remove shortcut %q: %w", filepath.Base(shortcutPath), err)
		}
	}

	if err := os.Remove(currentDir); err != nil {
		return fmt.Errorf("failed to unlink current %s: %w", currentDir, err)
	}

	return nil
}
