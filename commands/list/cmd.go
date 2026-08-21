package list

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

type Command struct{}

func (cmd *Command) Run() error {
	mochaDir := config.Current().MochaDirectory
	appDir := filepath.Join(mochaDir, "apps")

	apps, err := os.ReadDir(appDir)
	if err != nil {
		return fmt.Errorf("failed to read apps directory: %w", err)
	}

	headers := []string{"Name", "Version", "Bucket"}
	var rows [][]string

	for _, app := range apps {
		versions, err := os.ReadDir(filepath.Join(appDir, app.Name()))
		if err != nil {
			return fmt.Errorf("failed to read directory %q: %w", app.Name(), err)
		}

		for _, version := range versions {
			if version.Name() == "current" || !version.IsDir() {
				continue
			}

			info := manifest.Info{
				Name:         app.Name(),
				Bucket:       "",
				Version:      version.Name(),
				ManifestPath: "",
			}

			info, err := manifest.PopulateInfo(info, mochaDir)
			if err != nil {
				return fmt.Errorf("failed to fetch app details for %q: %w", app.Name(), err)
			}

			rows = append(rows, []string{info.Name, info.Version, info.Bucket})
		}
	}

	tableConfig := output.TableConfig{
		Spacing: 2,
		Alignments: []output.Alignment{
			output.LeftAlign,
			output.LeftAlign,
			output.LeftAlign,
		},
		BorderStyle: output.LightBorder,
	}

	if err := output.PrintTable(headers, rows, tableConfig); err != nil {
		return fmt.Errorf("failed to print app info table: %w", err)
	}
	return nil
}
