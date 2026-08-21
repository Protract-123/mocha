package cache

import (
	"github.com/Protract-123/mocha/commands/cache/clear"
	"github.com/Protract-123/mocha/commands/cache/show"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Show  *show.Command  `arg:"subcommand:show" help:"show all cache items"`
	Clear *clear.Command `arg:"subcommand:clear" help:"clear cache items"`
}

func (cmd *Command) Run() error {
	switch {
	case cmd.Show != nil:
		return cmd.Show.Run()
	case cmd.Clear != nil:
		return cmd.Clear.Run()
	default:
		return arg.ErrHelp
	}
}
