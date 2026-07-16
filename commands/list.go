package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

type ListCommand struct{}

func (cmd *ListCommand) Run() error {
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

			ref := manifest.Ref{
				Name:         app.Name(),
				Bucket:       "",
				Version:      version.Name(),
				ManifestPath: "",
			}

			ref, err := manifest.PopulateRef(ref, mochaDir)
			if err != nil {
				return fmt.Errorf("failed to fetch app details for %q: %w", app.Name(), err)
			}

			rows = append(rows, []string{ref.Name, ref.Version, ref.Bucket})
		}
	}

	rows = append(rows, []string{})

	if err := output.PrintTable(headers, rows); err != nil {
		return fmt.Errorf("failed to print app info table: %w", err)
	}
	return nil
}
