package commands

import (
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
		return err
	}

	if cmd.Add != nil {
		err := cmd.Add.Run(mochaDir)
		if err != nil {
			return err
		}
	}
	if cmd.Remove != nil {
		err := cmd.Remove.Run(mochaDir)
		if err != nil {
			return err
		}
	}
	if cmd.List != nil {
		err := cmd.List.Run(mochaDir)
		if err != nil {
			return err
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
		return err
	}

	if shimPath == "" {
		resolved, err := exec.LookPath(cmd.Path)
		if err != nil {
			return err
		}
		shimPath = resolved
	}

	err = shim.CreateShim(cmd.Name, shimPath, mochaDir)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *removeShimCommand) Run(mochaDir string) error {
	return shim.DeleteShim(cmd.Name, mochaDir)
}

func (cmd *listShimCommand) Run(mochaDir string) error {
	shims, err := shim.GetAllShims(mochaDir)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
