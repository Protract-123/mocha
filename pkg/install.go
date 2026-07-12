package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/shim"
)

type InstallInfo struct {
	Bucket  string `json:"bucket"`
	Version string `json:"version"`
	Arch    string `json:"architecture"`
}

type InstallOptions struct {
	SubDir string

	InnoSetup    bool
	RealFileName string
}

func InstallPackageFile(filePath string, installDir string, mochaDir string, options InstallOptions) error {
	extension := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filepath.Base(filePath), extension)

	tempDir := filepath.Join(mochaDir, "temp", fileName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", tempDir, err)
	}
	defer os.RemoveAll(filepath.Join(mochaDir, "temp"))

	switch extension {
	case ".zip":
		if err := fileops.ExtractZip(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	case ".msi":
		if err := fileops.ExtractMsi(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	case ".exe":
		if !options.InnoSetup {
			if options.RealFileName == "" {
				return fmt.Errorf("missing real file name for %s", filePath)
			}
			targetFilePath := filepath.Join(tempDir, options.RealFileName)

			if err := fileops.CopyFile(filePath, targetFilePath); err != nil {
				return fmt.Errorf("failed to move %s: %w", targetFilePath, err)
			}
		} else {
			if err := fileops.ExtractInnoSetup(filePath, tempDir, options.SubDir); err != nil {
				return fmt.Errorf("failed to extract InnoSetup from %s: %w", filePath, err)
			}
			options.SubDir = ""
		}
	default:
		if err := fileops.Extract7z(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	}

	extractedDir := filepath.Join(tempDir, options.SubDir)
	if err := fileops.MergeDir(extractedDir, installDir); err != nil {
		return fmt.Errorf("failed to merge %s into %s: %w", extractedDir, installDir, err)
	}

	return nil
}

func InstallApp(ref manifest.Ref, downloadArch string, force bool, mochaDir string) error {
	downloadResults, err := DownloadPackageFiles(ref, downloadArch, force, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to download manifest files: %w", err)
	}

	versionDir := filepath.Join(mochaDir, "apps", ref.Name, ref.Version)
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", versionDir, err)
	}

	innoSetup := manifest.GetManifestInnoSetup(ref.ManifestPath)
	for _, result := range downloadResults {
		installOptions := InstallOptions{
			SubDir:       result.Entry.SubDir,
			InnoSetup:    innoSetup,
			RealFileName: result.RealFilename,
		}

		if err := InstallPackageFile(result.DownloadPath, versionDir, mochaDir, installOptions); err != nil {
			return fmt.Errorf("failed to install %s: %w", result.Filename, err)
		}
	}

	currentDir := filepath.Join(mochaDir, "apps", ref.Name, "current")
	if err := os.Remove(currentDir); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove old junction %s: %w", currentDir, err)
	}

	if err := fileops.CreateJunction(versionDir, currentDir); err != nil {
		return fmt.Errorf("failed to create junction: %w", err)
	}

	binaries, err := manifest.GetManifestBin(ref.ManifestPath, downloadArch)
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

	shortcutEntries, err := manifest.GetManifestShortcut(ref.ManifestPath, downloadArch)
	if err != nil {
		return fmt.Errorf("failed to get shortcuts to create: %w", err)
	}

	shortcutDirectory := os.ExpandEnv(filepath.Join("$APPDATA", "Microsoft", "Windows", "Start Menu", "Programs", "Mocha Apps"))
	if err := os.MkdirAll(shortcutDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create shortcut directory: %w", err)
	}

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

	persistEntries, err := manifest.GetManifestPersist(ref.ManifestPath)
	if err != nil {
		return fmt.Errorf("failed to get persist entries: %w", err)
	}

	for _, persistEntry := range persistEntries {
		source := filepath.Join(mochaDir, "persist", ref.Name, persistEntry.Source)
		target := filepath.Join(currentDir, persistEntry.Target)

		var sourceExists bool
		var targetExists bool

		if _, err := os.Stat(source); err == nil {
			sourceExists = true
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to check if persist source exists: %w", err)
		}

		if _, err := os.Stat(target); err == nil {
			targetExists = true
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to check if persist target exists: %w", err)
		}

		switch {
		case !sourceExists && !targetExists:
			if err := os.MkdirAll(source, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create source directory %q: %w", source, err)
			}
		case !sourceExists && targetExists:
			if err := os.MkdirAll(filepath.Dir(source), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create parent directory for %q: %w", source, err)
			}
			if err := os.Rename(target, source); err != nil {
				return fmt.Errorf("failed to move target %q to %q: %w", target, source, err)
			}
		case sourceExists && targetExists:
			if err := os.Rename(target, target+".original"); err != nil {
				return fmt.Errorf("failed to move target %q to %q: %w", target, target+".original", err)
			}
		case sourceExists && !targetExists:
			break
		}

		info, err := os.Stat(source)
		if err != nil {
			return fmt.Errorf("failed to get persist source info: %w", err)
		}

		if info.IsDir() {
			if err := fileops.CreateJunction(source, target); err != nil {
				return fmt.Errorf("failed to symlink (junction) target to source: %w", err)
			}
		} else {
			if err := os.Link(source, target); err != nil {
				return fmt.Errorf("failed to symlink (hardlink) target to source: %w", err)
			}
		}
	}

	installJsonFile, err := os.Create(filepath.Join(currentDir, "install.json"))
	if err != nil {
		return fmt.Errorf("failed to create install.json: %w", err)
	}

	installInfo := InstallInfo{
		Bucket:  ref.Bucket,
		Version: ref.Version,
		Arch:    downloadArch,
	}

	jsonEncoder := json.NewEncoder(installJsonFile)
	if err := jsonEncoder.Encode(installInfo); err != nil {
		return fmt.Errorf("failed to write install info to install.json: %w", err)
	}

	return nil
}
