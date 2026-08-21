package setup

import (
	"fmt"
	"runtime"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/shim"
)

type Command struct {
	Language string `arg:"positional,required" help:"language of the shim binary"`
	Version  string `arg:"positional,required" help:"version of the shim binary"`
}

func (cmd *Command) Run() error {
	var selectedRelease shim.Release

	releases, err := shim.GetLatestShimReleases()
	if err != nil {
		return fmt.Errorf("failed to get latest shim releases: %w", err)
	}

	switch {
	case cmd.Language == "" && cmd.Version == "":
	// prompt user for release to download
	case cmd.Language != "" && cmd.Version == "":
		for _, entry := range releases {
			if entry.Language == cmd.Language {
				selectedRelease = entry
			}
		}
	case cmd.Language != "" && cmd.Version != "":
		for _, entry := range releases {
			if entry.Language == cmd.Language && entry.Version == cmd.Version {
				selectedRelease = entry
			}
		}
	case cmd.Language == "" && cmd.Version != "":
		return fmt.Errorf("version provided without a language specified")
	}

	if selectedRelease.Version == "" || selectedRelease.Language == "" {
		return fmt.Errorf("no valid release found with given parameters")
	}

	arch, err := getShimArch()
	if err != nil {
		return fmt.Errorf("invalid arch: %w", err)
	}

	if err := shim.InstallBinary(selectedRelease, arch, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to setup shim binary: %w", err)
	}

	return nil
}

func getShimArch() (string, error) {
	switch runtime.GOARCH {
	case "386":
		return "x86", nil
	case "amd64":
		return "x64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("cpu architecture %q is unsupported", runtime.GOARCH)
	}
}
