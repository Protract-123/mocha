package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

type InstallCommand struct {
	Apps       []string `arg:"positional,required" help:"apps to install (e.g. git, bat@0.26.1)"`
	Force      bool     `arg:"-f,--force" help:"ignore cache hits"`
	SkipVerify bool     `arg:"-s,--skip-verify" help:"skip hash check"`
}

func (cmd *InstallCommand) Run() error {
	downloadArch, err := pkg.DownloadArch()
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

		manifestJson, err := manifest.GetJson(manifestRef.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get manifest JSON: %w", err)
		}

		target := pkg.Package{
			Ref:  manifestRef,
			Json: manifestJson,
			Arch: downloadArch,
		}

		downloadOptions := pkg.DownloadOptions{
			Force:      cmd.Force,
			SkipVerify: cmd.SkipVerify,
		}

		downloadResults, err := pkg.Download(target, mochaDir, downloadOptions)
		if err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}

		if err := pkg.Install(target, downloadResults, mochaDir); err != nil {
			return fmt.Errorf("failed to install %q: %w", refString, err)
		}

		if err := pkg.Link(target.Ref, mochaDir); err != nil {
			return fmt.Errorf("failed to link app: %w", err)
		}

		output.LogSuccess("successfully installed %q", manifestRef.Name)
	}

	return nil
}
