package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/manifest"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	ManifestReferences []string `arg:"positional,required"`
	Force              bool     `arg:"-f,--force"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
	for _, refString := range cmd.ManifestReferences {
		if _, _, err := manifest.DownloadManifestFiles(refString, cmd.Force, mochaDir); err != nil {
			return fmt.Errorf("error downloading manifest files: %w", err)
		}
	}
	return nil
}
