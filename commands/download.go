package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	ManifestReferences []string `arg:"positional"`
	Force              bool     `arg:"-f,--force"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
	if cmd.ManifestReferences == nil {
		return errors.New("at least one manifest reference is required")
	}

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
			return err
		}

		downloadEntries, err := manifest.GetManifestDownloads(manifestRef.ManifestPath, downloadArch)
		if err != nil {
			return err
		}

		for _, entry := range downloadEntries {
			downloadPath := fileops.GetCachePath(mochaDir, manifestRef.Name, manifestRef.Version, entry.URL)
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
