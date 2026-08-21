package remove

import (
	"fmt"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
)

type Command struct {
	Name string `arg:"positional,required" help:"bucket name (e.g. main)"`
}

func (cmd *Command) Run() error {
	if err := bucket.DeleteBucket(cmd.Name, config.Current().MochaDirectory); err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", cmd.Name, err)
	}

	output.LogSuccess("successfully deleted bucket %q", cmd.Name)
	return nil
}
