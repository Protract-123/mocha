package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/pkg"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	Apps  []string `arg:"positional,required" help:"apps to download (e.g. git, bat@0.26.1)"`
	Force bool     `arg:"-f,--force" help:"ignore cache hits"`
}

func (cmd *DownloadCommand) Run(mochaDir string) error {
	downloadArch, err := pkg.GetDownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	for _, refString := range cmd.Apps {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", refString, err)
		}

		if _, err := pkg.DownloadPackageFiles(manifestRef, downloadArch, cmd.Force, mochaDir); err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}
	}
	return nil
}
