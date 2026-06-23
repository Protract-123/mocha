package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type ShimCommand struct {
	Add    *addShimCommand    `arg:"subcommand:add"`
	Remove *removeShimCommand `arg:"subcommand:remove"`
	List   *listShimCommand   `arg:"subcommand:list"`
	//Info   *infoShimCommand   `arg:"subcommand:info"`
	//Alter  *alterShimCommand  `arg:"subcommand:alter"`
}

type addShimCommand struct {
	Name string `arg:"positional,required"`
	Path string `arg:"positional,required"`
}
type removeShimCommand struct {
	Name string `arg:"positional,required"`
}
type listShimCommand struct{}

//type infoShimCommand struct{}
//type alterShimCommand struct{}

func (cmd *ShimCommand) Run(mochaDir string) error {
	err := shim.InitShims(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to init shims: %w", err)
	}

	if cmd.Add != nil {
		err := cmd.Add.Run(mochaDir)
		if err != nil {
			return fmt.Errorf("failed to add shims: %w", err)
		}
	}
	if cmd.Remove != nil {
		err := cmd.Remove.Run(mochaDir)
		if err != nil {
			return fmt.Errorf("failed to remove shims: %w", err)
		}
	}
	if cmd.List != nil {
		err := cmd.List.Run(mochaDir)
		if err != nil {
			return fmt.Errorf("failed to list shims: %w", err)
		}
	}

	return nil
}

func (cmd *addShimCommand) Run(mochaDir string) error {
	var shimPath string

	_, err := os.Stat(cmd.Path)
	if err == nil {
		shimPath = cmd.Path
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to confirm target's existence: %w", err)
	}

	if shimPath == "" {
		resolved, err := exec.LookPath(cmd.Path)
		if err != nil {
			return fmt.Errorf("failed to lookup target's path: %w", err)
		}
		shimPath = resolved
	}

	err = shim.CreateShim(cmd.Name, shimPath, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to create shim: %w", err)
	}

	return nil
}

func (cmd *removeShimCommand) Run(mochaDir string) error {
	err := shim.DeleteShim(cmd.Name, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to delete shim: %w", err)
	}

	return nil
}

func (cmd *listShimCommand) Run(mochaDir string) error {
	shims, err := shim.GetAllShims(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get all shims: %w", err)
	}

	headers := []string{"Name", "Destination"}
	rows := make([][]string, len(shims))

	for i, entry := range shims {
		rows[i] = []string{
			entry.Name,
			entry.Destination,
		}
	}

	err = output.PrintTable(headers, rows)
	if err != nil {
		return fmt.Errorf("failed to print shim info table: %w", err)
	}

	return nil
}
