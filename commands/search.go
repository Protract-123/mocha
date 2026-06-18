package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/output"
	"github.com/sahilm/fuzzy"
)

type SearchCommand struct {
	Query string `arg:"positional,required"`
	Count int    `default:"20" arg:"-c, --count"`
}

func (cmd *SearchCommand) Run(mochaDir string) error {
	cmd.Query = strings.ToLower(cmd.Query)
	bucketsDir := filepath.Join(mochaDir, "buckets")

	buckets, err := os.ReadDir(bucketsDir)
	if err != nil {
		return err
	}

	var allManifestNames []string
	var exactMatches []string

	for _, bucket := range buckets {
		if !bucket.IsDir() {
			continue
		}

		manifestDir := filepath.Join(bucketsDir, bucket.Name(), "bucket")
		manifests, err := os.ReadDir(manifestDir)
		if err != nil {
			return err
		}

		for _, manifest := range manifests {
			appName := strings.TrimSuffix(manifest.Name(), filepath.Ext(manifest.Name()))

			if appName == cmd.Query {
				exactMatches = append(exactMatches, appName)
				continue
			}

			allManifestNames = append(allManifestNames, appName)
		}
	}

	fuzzyResults := fuzzy.Find(cmd.Query, allManifestNames)

	if len(fuzzyResults) == 0 && len(exactMatches) == 0 {
		return fmt.Errorf("no results found")
	}

	limit := len(fuzzyResults)
	if cmd.Count > 0 && cmd.Count < limit {
		limit = cmd.Count
	}

	if len(exactMatches) != 0 {
		fmt.Println("\nExact Matches:")

		for _, result := range exactMatches {
			appDetails, err := app.ParseAppString(result)
			if err != nil {
				return err
			}

			appDetails, err = app.PopulateAppRef(appDetails, mochaDir)
			if err != nil {
				return err
			}

			fmt.Printf("%s - %s - %s\n", appDetails.Name, appDetails.Bucket, appDetails.Version)
		}

		fmt.Print("\n")
	}

	headers := []string{"Name", "Bucket", "Version"}
	rows := make([][]string, limit)

	for index, result := range fuzzyResults[:limit] {
		appDetails, err := app.ParseAppString(result.Str)
		if err != nil {
			return err
		}

		appDetails, err = app.PopulateAppRef(appDetails, mochaDir)
		if err != nil {
			return err
		}

		rows[index] = []string{appDetails.Name, appDetails.Bucket, appDetails.Version}
	}

	return output.PrintTable(headers, rows)
}
