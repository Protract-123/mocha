package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/fileops"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	AppReferences []string `arg:"positional"`
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

			err := fileops.DownloadFile(entry.URL, downloadPath)
			if err != nil {
				return fmt.Errorf("failed to download %s: %w", entry.URL, err)
			}

			err = fileops.VerifyHash(downloadPath, entry.Hash)
			if err != nil {
				err2 := os.Remove(downloadPath)
				if err2 != nil {
					return err2
				}

				return err
			}

			fmt.Printf("Downloaded and verified %s\n", downloadPath)
		}
	}

	return nil
}
