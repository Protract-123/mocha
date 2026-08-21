package shim

import (
	"github.com/Protract-123/mocha/commands/shim/add"
	"github.com/Protract-123/mocha/commands/shim/list"
	"github.com/Protract-123/mocha/commands/shim/remove"
	"github.com/Protract-123/mocha/commands/shim/setup"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Add    *add.Command    `arg:"subcommand:add" help:"add a shim"`
	Remove *remove.Command `arg:"subcommand:remove" help:"remove an existing shim"`
	List   *list.Command   `arg:"subcommand:list" help:"list added shims"`
	Setup  *setup.Command  `arg:"subcommand:setup" help:"initialize the shim binary"`
}

func (cmd *Command) Run() error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run()
	case cmd.Remove != nil:
		return cmd.Remove.Run()
	case cmd.List != nil:
		return cmd.List.Run()
	case cmd.Setup != nil:
		return cmd.Setup.Run()
	default:
		return arg.ErrHelp
	}
}
