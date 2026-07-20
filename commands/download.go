package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

// TODO: add more/better logging, like a progress bar

type DownloadCommand struct {
	Apps       []string `arg:"positional,required" help:"apps to download (e.g. git, bat@0.26.1)"`
	Force      bool     `arg:"-f,--force" help:"ignore cache hits"`
	SkipVerify bool     `arg:"-s,--skip-verify" help:"skip hash check"`
}

func (cmd *DownloadCommand) Run() error {
	downloadArch, err := pkg.DownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	mochaDir := config.Current().MochaDirectory

	for _, spec := range cmd.Apps {
		info, err := manifest.ParseSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse manifest spec %q: %w", spec, err)
		}

		info, err = manifest.PopulateInfo(info, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", spec, err)
		}

		manifestJson, err := manifest.GetJson(info.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get manifest JSON: %w", err)
		}

		target := pkg.Package{
			Info: info,
			Json: manifestJson,
			Arch: downloadArch,
		}

		options := pkg.DownloadOptions{
			Force:      cmd.Force,
			SkipVerify: cmd.SkipVerify,
		}

		if _, err := pkg.Download(target, mochaDir, options); err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}

		output.LogSuccess("successfully downloaded %q", spec)
	}
	return nil
}
