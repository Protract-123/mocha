package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	AppReferences []string `arg:"positional"`
	Force         bool     `arg:"-f,--force"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
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
		}
	}

	return nil
}
