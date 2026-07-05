package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/manifest"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	Apps  []string `arg:"positional,required" help:"apps to download (e.g. git, bat@0.26.1)"`
	Force bool     `arg:"-f,--force" help:"ignore cache hits"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
	for _, refString := range cmd.Apps {
		if _, _, err := manifest.DownloadManifestFiles(refString, cmd.Force, mochaDir); err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}
	}
	return nil
}
