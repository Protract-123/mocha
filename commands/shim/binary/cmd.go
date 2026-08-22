package binary

import (
	"github.com/Protract-123/mocha/commands/shim/binary/releases"
	"github.com/Protract-123/mocha/commands/shim/binary/setup"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Setup    *setup.Command    `arg:"subcommand:setup" help:"initialize the shim binary"`
	Releases *releases.Command `arg:"subcommand:releases" help:"list shim releases"`
}

func (cmd *Command) Run() error {
	switch {
	case cmd.Releases != nil:
		return cmd.Releases.Run()
	case cmd.Setup != nil:
		return cmd.Setup.Run()
	default:
		return arg.ErrHelp
	}
}
