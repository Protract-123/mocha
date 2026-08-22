package add

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type Command struct {
	Name string `arg:"positional,required" help:"name for the shim (e.g. bat)"`
	Path string `arg:"positional,required" help:"path to the executable, or a name resolvable on PATH (e.g. C:/path/to/bat.exe or bat)"`
}

func (cmd *Command) Run() error {
	shimPath := cmd.Path

	mochaDir := config.Current().MochaDirectory

	if _, err := os.Stat(filepath.Join(mochaDir, "shim.exe")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("shim.exe does not exist, setup shim binary with shim binary setup")
	} else if err != nil {
		return fmt.Errorf("failed to confirm target's existence: %w", err)
	}

	if _, err := os.Stat(cmd.Path); errors.Is(err, os.ErrNotExist) {
		resolved, err := exec.LookPath(cmd.Path)
		if err != nil {
			return fmt.Errorf("failed to lookup target's path: %w", err)
		}
		shimPath = resolved
	} else if err != nil {
		return fmt.Errorf("failed to confirm target's existence: %w", err)
	}

	if err := shim.CreateShim(cmd.Name, shimPath, mochaDir); err != nil {
		return fmt.Errorf("failed to create shim: %w", err)
	}

	output.LogSuccess("successfully created shim %q", cmd.Name)
	return nil
}
