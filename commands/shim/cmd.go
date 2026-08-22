package shim

import (
	"github.com/Protract-123/mocha/commands/shim/add"
	"github.com/Protract-123/mocha/commands/shim/binary"
	"github.com/Protract-123/mocha/commands/shim/list"
	"github.com/Protract-123/mocha/commands/shim/remove"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Add    *add.Command    `arg:"subcommand:add" help:"add a shim"`
	Remove *remove.Command `arg:"subcommand:remove" help:"remove an existing shim"`
	List   *list.Command   `arg:"subcommand:list" help:"list added shims"`
	Binary *binary.Command `arg:"subcommand:binary" help:"manage the shim binary"`
}

func (cmd *Command) Run() error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run()
	case cmd.Remove != nil:
		return cmd.Remove.Run()
	case cmd.List != nil:
		return cmd.List.Run()
	case cmd.Binary != nil:
		return cmd.Binary.Run()
	default:
		return arg.ErrHelp
	}
}
