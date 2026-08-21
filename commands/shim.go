package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
	"github.com/alexflint/go-arg"
)

type ShimCommand struct {
	Add    *addShimCommand    `arg:"subcommand:add" help:"add a shim"`
	Remove *removeShimCommand `arg:"subcommand:remove" help:"remove an existing shim"`
	List   *listShimCommand   `arg:"subcommand:list" help:"list added shims"`
	Init   *initShimCommand   `arg:"subcommand:init" help:"initialize the shim binary"`
}

type addShimCommand struct {
	Name string `arg:"positional,required" help:"name for the shim (e.g. bat)"`
	Path string `arg:"positional,required" help:"path to the executable, or a name resolvable on PATH (e.g. C:/path/to/bat.exe or bat)"`
}
type removeShimCommand struct {
	Name string `arg:"positional,required" help:"name of the shim to remove"`
}
type listShimCommand struct{}
type initShimCommand struct {
	Language string `arg:"positional,required" help:"language of the shim binary"`
	Version  string `arg:"positional,required" help:"version of the shim binary"`
}

func (cmd *ShimCommand) Run() error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run()
	case cmd.Remove != nil:
		return cmd.Remove.Run()
	case cmd.List != nil:
		return cmd.List.Run()
	case cmd.Init != nil:
		return cmd.Init.Run()
	default:
		return arg.ErrHelp
	}
}

func (cmd *addShimCommand) Run() error {
	shimPath := cmd.Path

	if _, err := os.Stat(cmd.Path); errors.Is(err, os.ErrNotExist) {
		resolved, err := exec.LookPath(cmd.Path)
		if err != nil {
			return fmt.Errorf("failed to lookup target's path: %w", err)
		}
		shimPath = resolved
	} else if err != nil {
		return fmt.Errorf("failed to confirm target's existence: %w", err)
	}

	if err := shim.CreateShim(cmd.Name, shimPath, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to create shim: %w", err)
	}

	output.LogSuccess("successfully created shim %q", cmd.Name)
	return nil
}

func (cmd *removeShimCommand) Run() error {
	if err := shim.DeleteShim(cmd.Name, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to delete shim: %w", err)
	}

	output.LogSuccess("successfully deleted shim %q", cmd.Name)
	return nil
}

func (cmd *listShimCommand) Run() error {
	shims, err := shim.GetAllShims(config.Current().MochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to get all shims: %w", err)
	}

	if len(shims) == 0 {
		return fmt.Errorf("no shims found")
	}

	headers := []string{"Name", "Destination"}
	rows := make([][]string, len(shims))

	for i, entry := range shims {
		rows[i] = []string{
			entry.Name,
			entry.Target,
		}
	}

	tableConfig := output.TableConfig{
		Spacing: 2,
		Alignments: []output.Alignment{
			output.LeftAlign,
			output.LeftAlign,
		},
		BorderStyle: output.LightBorder,
	}

	if err := output.PrintTable(headers, rows, tableConfig); err != nil {
		return fmt.Errorf("failed to print shim info table: %w", err)
	}

	return nil
}

func (cmd *initShimCommand) Run() error {
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
		return fmt.Errorf("failed to init shim binary: %w", err)
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