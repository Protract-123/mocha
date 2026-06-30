package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type InstallCommand struct {
	ManifestReferences []string `arg:"positional,required"`
	Force              bool     `arg:"-f,--force"`
}

func (cmd *InstallCommand) Run(mochaDir string) error {
	for _, refString := range cmd.ManifestReferences {
		manifestRef, downloadResults, err := manifest.DownloadManifestFiles(refString, cmd.Force, mochaDir)
		if err != nil {
			return fmt.Errorf("error downloading manifest files: %w", err)
		}

		versionDir := filepath.Join(mochaDir, "apps", manifestRef.Name, manifestRef.Version)
		if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", versionDir, err)
		}

		innoSetup := manifest.GetManifestInnoSetup(manifestRef.ManifestPath)

		for _, result := range downloadResults {
			installOptions := manifest.InstallOptions{
				SubDir:       result.Entry.SubDir,
				InnoSetup:    innoSetup,
				RealFileName: result.RealFilename,
			}

			if err := manifest.InstallManifestFile(result.DownloadPath, versionDir, mochaDir, installOptions); err != nil {
				return fmt.Errorf("failed to install %s: %w", result.Filename, err)
			}
		}

		currentDir := filepath.Join(mochaDir, "apps", manifestRef.Name, "current")
		if err := os.Remove(currentDir); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to remove old junction %s: %w", currentDir, err)
		}

		if err := fileops.CreateJunction(versionDir, currentDir); err != nil {
			return fmt.Errorf("failed to create junction: %w", err)
		}

		binaries, err := manifest.GetManifestBin(manifestRef.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get binaries to shim: %w", err)
		}

		for _, binary := range binaries {
			shimName := strings.TrimSuffix(filepath.Base(binary), filepath.Ext(binary))
			shimPath := filepath.Join(currentDir, binary)
			if err := shim.CreateShim(shimName, shimPath, mochaDir); err != nil {
				return fmt.Errorf("failed to create shim %s: %w", shimName, err)
			}
		}

		output.LogOutput(fmt.Sprintf("Installed %s", manifestRef.Name))

	}

	return nil
}
