package list

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
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
		if !app.IsDir() {
			continue
		}

		versions, err := os.ReadDir(filepath.Join(appDir, app.Name()))
		if err != nil {
			return fmt.Errorf("failed to read directory %q: %w", app.Name(), err)
		}

		for _, version := range versions {
			if version.Name() == "current" || !version.IsDir() {
				continue
			}

			installJson, err := os.ReadFile(filepath.Join(appDir, app.Name(), version.Name(), "install.json"))
			if err != nil {
				return fmt.Errorf("failed to read install JSON of %q: %w", app.Name(), err)
			}

			installInfo := pkg.InstallInfo{}
			if err := json.Unmarshal(installJson, &installInfo); err != nil {
				return fmt.Errorf("failed to unmarshal install JSON: %w", err)
			}

			rows = append(rows, []string{app.Name(), version.Name(), installInfo.Bucket})
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
