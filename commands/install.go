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
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", refString, err)
		}

		downloadArch, err := manifest.GetSystemArch()
		if err != nil {
			return fmt.Errorf("failed to get system architecture: %w", err)
		}

		downloadEntries, err := manifest.GetManifestDownloads(manifestRef.ManifestPath, downloadArch)
		if err != nil {
			return fmt.Errorf("failed to get manifest downloads: %w", err)
		}

		var downloadPaths []string

		for _, entry := range downloadEntries {
			downloadPath := fileops.GetCachePath(mochaDir, manifestRef.Name, manifestRef.Version, entry.URL)
			filename := filepath.Base(downloadPath)

			if _, err := os.Stat(downloadPath); err != nil || cmd.Force {
				output.LogOutput(fmt.Sprintf("Downloading %s to %s", entry.URL, downloadPath))
				if err := fileops.DownloadFile(entry.URL, downloadPath); err != nil {
					return fmt.Errorf("failed to download %s: %w", filename, err)
				}
				output.LogOutput(fmt.Sprintf("Downloaded %s", filename))
			} else {
				output.LogOutput(fmt.Sprintf("Cache hit, skipping %s", filename))
			}

			if err := fileops.VerifyHash(downloadPath, entry.Hash); err != nil {
				_ = os.Remove(downloadPath)
				return fmt.Errorf("failed to verify %s: %w", filename, err)
			}

			output.LogOutput(fmt.Sprintf("Verified %s\n", filename))
			downloadPaths = append(downloadPaths, downloadPath)
		}

		appDir := filepath.Join(mochaDir, "apps", manifestRef.Name)
		versionDir := filepath.Join(appDir, manifestRef.Version)
		currentDir := filepath.Join(appDir, "current")

		if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", versionDir, err)
		}

		for _, downloadPath := range downloadPaths {
			if err := fileops.ExtractFile(downloadPath, versionDir); err != nil {
				return fmt.Errorf("failed to extract %s: %w", filepath.Base(downloadPath), err)
			}
		}

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
