package bucket

import (
	"github.com/Protract-123/mocha/commands/bucket/add"
	"github.com/Protract-123/mocha/commands/bucket/known"
	"github.com/Protract-123/mocha/commands/bucket/list"
	"github.com/Protract-123/mocha/commands/bucket/remove"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Add    *add.Command    `arg:"subcommand:add" help:"add a bucket by name or git repository URL"`
	Remove *remove.Command `arg:"subcommand:remove" help:"remove an installed bucket"`
	List   *list.Command   `arg:"subcommand:list" help:"list installed buckets"`
	Known  *known.Command  `arg:"subcommand:known" help:"list known buckets available to add by name"`
}

func (cmd *Command) Run() error {
	switch {
	case cmd.Add != nil:
		return cmd.Add.Run()
	case cmd.Remove != nil:
		return cmd.Remove.Run()
	case cmd.List != nil:
		return cmd.List.Run()
	case cmd.Known != nil:
		return cmd.Known.Run()
	default:
		return arg.ErrHelp
	}
}
