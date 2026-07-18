package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
	"github.com/alexflint/go-arg"
)

type ShimCommand struct {
	Add    *addShimCommand    `arg:"subcommand:add" help:"add a shim"`
	Remove *removeShimCommand `arg:"subcommand:remove" help:"remove an existing shim"`
	List   *listShimCommand   `arg:"subcommand:list" help:"list added shims"`
}

type addShimCommand struct {
	Name string `arg:"positional,required" help:"name for the shim (e.g. bat)"`
	Path string `arg:"positional,required" help:"path to the executable, or a name resolvable on PATH (e.g. C:/path/to/bat.exe or bat)"`
}
type removeShimCommand struct {
	Name string `arg:"positional,required" help:"name of the shim to remove"`
}
type listShimCommand struct{}

func (cmd *ShimCommand) Run() error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run()
	case cmd.Remove != nil:
		return cmd.Remove.Run()
	case cmd.List != nil:
		return cmd.List.Run()
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

	return nil
}

func (cmd *removeShimCommand) Run() error {
	if err := shim.DeleteShim(cmd.Name, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to delete shim: %w", err)
	}

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
