package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

type InstallCommand struct {
	Apps  []string `arg:"positional,required" help:"apps to install (e.g. git, bat@0.26.1)"`
	Force bool     `arg:"-f,--force" help:"ignore cache hits"`
}

func (cmd *InstallCommand) Run() error {
	downloadArch, err := pkg.GetDownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	mochaDir := config.Current().MochaDirectory

	for _, refString := range cmd.Apps {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", refString, err)
		}

		if err := pkg.InstallApp(manifestRef, downloadArch, cmd.Force, mochaDir); err != nil {
			return fmt.Errorf("failed to install %q: %w", refString, err)
		}

		output.LogOutput(fmt.Sprintf("Installed %s", manifestRef.Name))
	}

	return nil
}
