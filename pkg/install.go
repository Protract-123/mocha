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
)

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

	manifestJson, err := manifest.GetJson(ref.ManifestPath)
	if err != nil {
		return fmt.Errorf("failed to get manifest JSON: %w", err)
	}

	innoSetup := manifest.GetInnoSetup(manifestJson)
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

	if err := createPersistLinks(ref, manifestJson, mochaDir); err != nil {
		return fmt.Errorf("failed to create persist links: %w", err)
	}

	installJsonFile, err := os.Create(filepath.Join(versionDir, "install.json"))
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

	if err := fileops.CopyFile(ref.ManifestPath, filepath.Join(versionDir, "manifest.json")); err != nil {
		return fmt.Errorf("failed to copy app manifest to install location: %w", err)
	}

	if err := LinkApp(ref, mochaDir); err != nil {
		return fmt.Errorf("failed to link app: %w", err)
	}

	return nil
}

func createPersistLinks(ref manifest.Ref, manifestJson map[string]any, mochaDir string) error {
	versionDir := filepath.Join(mochaDir, "apps", ref.Name, ref.Version)
	persistEntries := manifest.GetPersistEntries(manifestJson)

	for _, persistEntry := range persistEntries {
		source := filepath.Join(mochaDir, "persist", ref.Name, persistEntry.Source)
		target := filepath.Join(versionDir, persistEntry.Target)

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

	return nil
}
