package remove

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type Command struct {
	Name string `arg:"positional,required" help:"name of the shim to remove"`
}

func (cmd *Command) Run() error {
	if err := shim.DeleteShim(cmd.Name, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to delete shim: %w", err)
	}

	output.LogSuccess("successfully deleted shim %q", cmd.Name)
	return nil
}
