package known

import (
	"fmt"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
)

type Command struct{}

func (cmd *Command) Run() error {
	knownBuckets, err := bucket.GetKnownBuckets(config.Current().MochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to get known buckets: %w", err)
	}

	headers := []string{"Name", "Source"}
	rows := make([][]string, len(knownBuckets))

	for i, entry := range knownBuckets {
		rows[i] = []string{entry.Name, entry.Source}
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
		return fmt.Errorf("failed to display bucket metadata: %w", err)
	}

	return nil
}
