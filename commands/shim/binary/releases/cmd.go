package releases

import (
	"fmt"

	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/shim"
)

type Command struct {
	All bool `arg:"-a,--all" help:"list all releases, including old ones"`
}

func (cmd *Command) Run() error {
	var releases []shim.Release
	var err error

	if cmd.All {
		releases, err = shim.GetShimReleases()
	} else {
		releases, err = shim.GetLatestShimReleases()
	}

	if err != nil {
		return fmt.Errorf("failed to fetch shim releases: %w", err)
	}

	headers := []string{"Language", "Version", "Published"}
	rows := make([][]string, len(releases))

	for index, release := range releases {
		rows[index] = []string{
			release.Language,
			release.Version,
			release.PublishedAt.Format("02-01-2006"),
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
