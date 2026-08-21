package list

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type Command struct{}

func (cmd *Command) Run() error {
	shims, err := shim.GetAllShims(config.Current().MochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to get all shims: %w", err)
	}

	if len(shims) == 0 {
		return fmt.Errorf("no shims found")
	}

	headers := []string{"Name", "Destination"}
	rows := make([][]string, len(shims))

	for i, entry := range shims {
		rows[i] = []string{
			entry.Name,
			entry.Target,
		}
	}

	tableConfig := output.TableConfig{
		Spacing: 2,
		Alignments: []output.Alignment{
			output.LeftAlign,
			output.LeftAlign,
		},
		BorderStyle: output.LightBorder,
	}

	if err := output.PrintTable(headers, rows, tableConfig); err != nil {
		return fmt.Errorf("failed to print shim info table: %w", err)
	}

	return nil
}
