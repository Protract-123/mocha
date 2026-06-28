package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	ManifestReferences []string `arg:"positional,required"`
	Force              bool     `arg:"-f,--force"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
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
		}
	}

	return nil
}
