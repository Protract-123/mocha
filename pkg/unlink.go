package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/shim"
)

// TODO: This works on wrong info if a new version removes stuff which existed earlier

func UnlinkApp(ref manifest.Ref, downloadArch string, mochaDir string, versionDir string) error {
	manifestJson, err := manifest.GetJson(ref.ManifestPath)
	if err != nil {
		return fmt.Errorf("failed to get manifest JSON: %w", err)
	}

	binaries, err := manifest.GetExecutableEntries(manifestJson, downloadArch)
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

	shortcutEntries := manifest.GetShortcutEntries(manifestJson, downloadArch)
	for _, shortcutEntry := range shortcutEntries {
		shortcutPath := filepath.Join(shortcutDirectory, fmt.Sprintf("%s.lnk", shortcutEntry.Name))

		if err := os.Remove(shortcutPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to remove shortcut %q: %w", filepath.Base(shortcutPath), err)
		}
	}

	currentDir := filepath.Join(mochaDir, "apps", ref.Name, "current")

	currentTarget, err := os.Readlink(currentDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to read current link %s: %w", currentDir, err)
	}

	absVersionDir, err := filepath.Abs(versionDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", versionDir, err)
	}

	if filepath.Clean(currentTarget) != filepath.Clean(absVersionDir) {
		return nil
	}

	if err := os.Remove(currentDir); err != nil {
		return fmt.Errorf("failed to unlink current %s: %w", currentDir, err)
	}

	return nil
}
