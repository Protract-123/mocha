package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

type OutdatedCommand struct{}

func (cmd *OutdatedCommand) Run() error {
	mochaDir := config.Current().MochaDirectory
	appDir := filepath.Join(mochaDir, "apps")

	apps, err := os.ReadDir(appDir)
	if err != nil {
		return fmt.Errorf("failed to read apps directory: %w", err)
	}

	headers := []string{"Name", "Current Version", "New Version"}
	var rows [][]string

	for _, app := range apps {
		if !app.IsDir() {
			continue
		}

		installedVersion, err := pkg.GetActiveVersion(app.Name(), mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get current version of %q: %w", app.Name(), err)
		}

		if installedVersion == "" {
			continue
		}

		installJson, err := os.ReadFile(filepath.Join(appDir, app.Name(), installedVersion, "install.json"))
		if err != nil {
			return fmt.Errorf("failed to read install JSON of %q: %w", app.Name(), err)
		}

		installInfo := pkg.InstallInfo{}
		if err := json.Unmarshal(installJson, &installInfo); err != nil {
			return fmt.Errorf("failed to unmarshal install JSON: %w", err)
		}

		info := manifest.Info{
			Name:         app.Name(),
			Bucket:       installInfo.Bucket,
			Version:      "",
			ManifestPath: "",
		}

		info, err = manifest.PopulateInfo(info, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get latest version of %q: %w", app.Name(), err)
		}

		if manifest.CompareVersions(installedVersion, info.Version) == 1 {
			rows = append(rows, []string{info.Name, installedVersion, info.Version})
		}
	}

	if len(rows) <= 0 {
		return fmt.Errorf("no apps are outdated")
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
