package list

import (
	"fmt"
	"strconv"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
)

type Command struct{}

func (cmd *Command) Run() error {
	mochaDir := config.Current().MochaDirectory

	bucketMetadata, err := bucket.GetAllBucketMetadata(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get bucket metadata: %w", err)
	}

	if len(bucketMetadata) == 0 {
		return fmt.Errorf("no buckets to show")
	}

	headers := []string{"Name", "Source", "Updated", "Manifests"}
	rows := make([][]string, len(bucketMetadata))

	for index, entry := range bucketMetadata {
		rows[index] = []string{
			entry.Name,
			entry.Source,
			entry.LastUpdated.Format("02-01-2006 15:04:05"),
			strconv.Itoa(entry.ManifestCount),
		}
	}

	tableConfig := output.TableConfig{
		Spacing: 2,
		Alignments: []output.Alignment{
			output.LeftAlign,
			output.LeftAlign,
			output.LeftAlign,
			output.RightAlign,
		},
		BorderStyle: output.LightBorder,
	}

	if err := output.PrintTable(headers, rows, tableConfig); err != nil {
		return fmt.Errorf("failed to display bucket metadata: %w", err)
	}

	return nil
}
