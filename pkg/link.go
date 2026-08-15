package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/shim"
)

// TODO: add env_set and env_add_path support

func Link(info manifest.Info, mochaDir string) error {
	if info.Name == "" {
		return fmt.Errorf("package info doesn't have a name")
	}
	if info.Version == "" {
		return fmt.Errorf("package info doesn't have a version")
	}

	if currentVersion, err := GetActiveVersion(info.Name, mochaDir); err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	} else if currentVersion != "" {
		return fmt.Errorf("package %s is already linked", info.Name)
	}

	versionDir := filepath.Join(mochaDir, "apps", info.Name, info.Version)
	currentDir := filepath.Join(mochaDir, "apps", info.Name, "current")

	installInfo, err := GetInstallInfo(versionDir)
	if err != nil {
		return fmt.Errorf("failed to get app install info: %w", err)
	}

	if err := fileops.CreateJunction(versionDir, currentDir); err != nil {
		return fmt.Errorf("failed to create junction: %w", err)
	}

	binaries, err := manifest.GetExecutableEntries(installInfo.ManifestJson, installInfo.Arch)
	if err != nil {
		return fmt.Errorf("failed to get binaries to shim: %w", err)
	}

	for _, binary := range binaries {
		shimName := strings.TrimSuffix(filepath.Base(binary.Alias), filepath.Ext(binary.Alias))
		shimPath := filepath.Join(currentDir, binary.Exe)
		if err := shim.CreateShim(shimName, shimPath, mochaDir); err != nil {
			return fmt.Errorf("failed to create shim %s: %w", shimName, err)
		}
	}

	shortcutDirectory := os.ExpandEnv(filepath.Join("$APPDATA", "Microsoft", "Windows", "Start Menu", "Programs", "Mocha Apps"))
	if err := os.MkdirAll(shortcutDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create shortcut directory: %w", err)
	}

	shortcutEntries := manifest.GetShortcutEntries(installInfo.ManifestJson, installInfo.Arch)
	for _, shortcutEntry := range shortcutEntries {
		shortcutPath := filepath.Join(shortcutDirectory, fmt.Sprintf("%s.lnk", shortcutEntry.Name))
		shortcutName := filepath.Base(shortcutPath)

		if err := os.MkdirAll(filepath.Dir(shortcutPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create shortcut directory: %w", err)
		}

		shortcut := fileops.Shortcut{
			ShortcutPath:     shortcutPath,
			Target:           filepath.Join(currentDir, shortcutEntry.Exe),
			WorkingDirectory: currentDir,
			Arguments:        shortcutEntry.Args,
		}

		if shortcutEntry.Icon != "" {
			shortcut.IconLocation = filepath.Join(currentDir, shortcutEntry.Icon)
		} else {
			shortcut.IconLocation = shortcut.Target
		}

		if err := fileops.CreateShortcut(shortcut); err != nil {
			return fmt.Errorf("failed to create shortcut %q: %w", shortcutName, err)
		}
	}

	return nil
}

func Unlink(appName string, mochaDir string) error {
	currentDir := filepath.Join(mochaDir, "apps", appName, "current")

	installInfo, err := GetInstallInfo(currentDir)
	if err != nil {
		return fmt.Errorf("failed to get app install info: %w", err)
	}

	binaries, err := manifest.GetExecutableEntries(installInfo.ManifestJson, installInfo.Arch)
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

	shortcutEntries := manifest.GetShortcutEntries(installInfo.ManifestJson, installInfo.Arch)
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
