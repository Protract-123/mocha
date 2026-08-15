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

// TODO: env_set requires variable parsing (e.g. $dir)

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

	envEntries := manifest.GetEnvEntries(installInfo.ManifestJson, installInfo.Arch)
	for key, value := range envEntries {
		if err := fileops.SetEnvironmentVariable(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %q: %w", key, err)
		}
	}

	pathEntries := manifest.GetPathEntries(installInfo.ManifestJson, installInfo.Arch)
	if len(pathEntries) > 0 {
		currentPath, err := fileops.GetEnvironmentVariable("Path")
		if err != nil {
			return fmt.Errorf("failed to get Path environment variable: %w", err)
		}

		existing := make(map[string]bool)
		for entry := range strings.SplitSeq(currentPath, ";") {
			if entry != "" {
				existing[strings.ToLower(entry)] = true
			}
		}

		additions := make([]string, 0, len(pathEntries))
		for _, pathEntry := range pathEntries {
			pathAddition := filepath.Join(currentDir, pathEntry)
			key := strings.ToLower(pathAddition)
			if existing[key] {
				continue
			}
			existing[key] = true
			additions = append(additions, pathAddition)
		}

		if len(additions) > 0 {
			newPath := strings.Join(additions, ";")
			if currentPath != "" {
				newPath += ";" + currentPath
			}
			if err := fileops.SetEnvironmentVariable("Path", newPath); err != nil {
				return fmt.Errorf("failed to update Path environment variable: %w", err)
			}
		}
	}

	if err := fileops.PropagateEnvironment(); err != nil {
		return fmt.Errorf("failed to propagate environment changes: %w", err)
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

	envEntries := manifest.GetEnvEntries(installInfo.ManifestJson, installInfo.Arch)
	for key := range envEntries {
		if err := fileops.RemoveEnvironmentVariable(key); err != nil {
			return fmt.Errorf("failed to remove environment variable %q: %w", key, err)
		}
	}

	pathEntries := manifest.GetPathEntries(installInfo.ManifestJson, installInfo.Arch)
	if len(pathEntries) > 0 {
		currentPath, err := fileops.GetEnvironmentVariable("Path")
		if err != nil {
			return fmt.Errorf("failed to get Path environment variable: %w", err)
		}

		pathRemovals := make(map[string]bool, len(pathEntries))
		for _, pathEntry := range pathEntries {
			pathRemovals[strings.ToLower(filepath.Join(currentDir, pathEntry))] = true
		}

		remaining := make([]string, 0, len(pathEntries))
		for currentPathEntry := range strings.SplitSeq(currentPath, ";") {
			if currentPathEntry == "" || pathRemovals[strings.ToLower(currentPathEntry)] {
				continue
			}
			remaining = append(remaining, currentPathEntry)
		}

		newPath := strings.Join(remaining, ";")
		if err := fileops.SetEnvironmentVariable("Path", newPath); err != nil {
			return fmt.Errorf("failed to update Path environment variable: %w", err)
		}
	}

	if err := fileops.PropagateEnvironment(); err != nil {
		return fmt.Errorf("failed to propagate environment changes: %w", err)
	}

	if err := os.Remove(currentDir); err != nil {
		return fmt.Errorf("failed to unlink current %s: %w", currentDir, err)
	}

	return nil
}
