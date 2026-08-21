package install

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

type Command struct {
	Apps       []string `arg:"positional,required" help:"apps to install (e.g. git, bat@0.26.1)"`
	Force      bool     `arg:"-f,--force" help:"ignore cache hits"`
	SkipVerify bool     `arg:"-s,--skip-verify" help:"skip hash check"`
}

// TODO: error if package already installed

func (cmd *Command) Run() error {
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

		downloadOptions := pkg.DownloadOptions{
			Force:      cmd.Force,
			SkipVerify: cmd.SkipVerify,
		}

		downloadResults, err := pkg.Download(target, mochaDir, downloadOptions)
		if err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}

		if err := pkg.Install(target, downloadResults, mochaDir); err != nil {
			return fmt.Errorf("failed to install %q: %w", spec, err)
		}

		if err := pkg.Link(target.Info, mochaDir); err != nil {
			return fmt.Errorf("failed to link app: %w", err)
		}

		output.LogSuccess("successfully installed %q", info.Name)
	}

	return nil
}
