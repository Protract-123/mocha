package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type InstallCommand struct {
	AppReferences []string `arg:"positional"`
	Force         bool     `arg:"-f,--force"`
}

func (cmd InstallCommand) Run(mochaDir string) error {
	if cmd.AppReferences == nil {
		return errors.New("at least one app reference is required")
	}

	for _, appString := range cmd.AppReferences {
		appRef, err := app.ParseAppString(appString)
		if err != nil {
			return fmt.Errorf("failed to parse app ref %q: %w", appString, err)
		}

		appRef, err = app.PopulateAppRef(appRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", appString, err)
		}

		downloadArch, err := app.GetSystemArch()
		if err != nil {
			return err
		}

		downloadEntries, err := app.GetManifestDownloads(appRef.ManifestPath, downloadArch)
		if err != nil {
			return err
		}

		var downloadPaths []string

		for _, entry := range downloadEntries {
			downloadPath := fileops.GetCachePath(mochaDir, appRef.Name, appRef.Version, entry.URL)
			filename := filepath.Base(downloadPath)

			_, err = os.Stat(downloadPath)
			if err != nil || cmd.Force {
				output.LogOutput(fmt.Sprintf("Downloading %s to %s", entry.URL, downloadPath))
				err = fileops.DownloadFile(entry.URL, downloadPath)
				if err != nil {
					return fmt.Errorf("failed to download %s: %w", filename, err)
				}
				output.LogOutput(fmt.Sprintf("Downloaded %s", filename))
			} else {
				output.LogOutput(fmt.Sprintf("Cache hit, skipping %s", filename))
			}

			err = fileops.VerifyHash(downloadPath, entry.Hash)
			if err != nil {
				err2 := os.Remove(downloadPath)
				if err2 != nil {
					return fmt.Errorf("failed to remove %s after invalid hash: %w", filename, err2)
				}

				return fmt.Errorf("failed to verify %s: %w", filename, err)
			}

			output.LogOutput(fmt.Sprintf("Verified %s\n", filename))
			downloadPaths = append(downloadPaths, downloadPath)
		}

		appDir := filepath.Join(mochaDir, "apps", appRef.Name)
		versionDir := filepath.Join(appDir, appRef.Version)
		currentDir := filepath.Join(appDir, "current")

		err = os.MkdirAll(versionDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", versionDir, err)
		}

		for _, downloadPath := range downloadPaths {
			err = fileops.ExtractFile(downloadPath, versionDir)
			if err != nil {
				return fmt.Errorf("failed to extract %s: %w", filepath.Base(downloadPath), err)
			}
		}

		if _, err := os.Stat(currentDir); err == nil {
			err := os.Remove(currentDir)
			if err != nil {
				return fmt.Errorf("failed to remove old junction %s: %w", currentDir, err)
			}
		}

		err = fileops.CreateJunction(versionDir, currentDir)
		if err != nil {
			return fmt.Errorf("failed to create junction: %w", err)
		}

		binaries, err := app.GetManifestBin(appRef.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get binaries to shim: %w", err)
		}

		for _, binary := range binaries {
			shimName := strings.TrimSuffix(filepath.Base(binary), filepath.Ext(binary))
			shimPath := filepath.Join(currentDir, binary)
			err := shim.CreateShim(shimName, shimPath, mochaDir)
			if err != nil {
				return fmt.Errorf("failed to create shim %s: %w", shimName, err)
			}
		}

		output.LogOutput(fmt.Sprintf("Installed %s", appRef.Name))
	}

	return nil
}
