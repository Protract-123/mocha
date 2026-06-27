package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/manifest"
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
		return fmt.Errorf("failed to read buckets directory: %w", err)
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
			return fmt.Errorf("failed to read %s's manifest directory: %w", bucket.Name(), err)
		}

		for _, entry := range manifests {
			appName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

			if appName == cmd.Query {
				exactMatches = append(exactMatches, appName)
				continue
			}

			allManifestNames = append(allManifestNames, appName)
		}
	}

	fuzzyResults := fuzzy.Find(cmd.Query, allManifestNames)

	if len(fuzzyResults) == 0 && len(exactMatches) == 0 {
		return fmt.Errorf("no results found in buckets")
	}

	limit := len(fuzzyResults)
	if cmd.Count > 0 && cmd.Count < limit {
		limit = cmd.Count
	}

	if len(exactMatches) != 0 {
		fmt.Println("\nExact Matches:")

		for _, result := range exactMatches {
			manifestRef, err := manifest.ParseRefString(result)
			if err != nil {
				return fmt.Errorf("failed to parse result %s: %w", result, err)
			}

			manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
			if err != nil {
				return fmt.Errorf("failed to get %s manifest details: %w", result, err)
			}

			fmt.Printf("%s - %s - %s\n", manifestRef.Name, manifestRef.Bucket, manifestRef.Version)
		}

		fmt.Print("\n")
	}

	headers := []string{"Name", "Bucket", "Version"}
	rows := make([][]string, limit)

	for index, result := range fuzzyResults[:limit] {
		appDetails, err := manifest.ParseRefString(result.Str)
		if err != nil {
			return fmt.Errorf("failed to parse result %s: %w", result.Str, err)
		}

		appDetails, err = manifest.PopulateRef(appDetails, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %s manifest details: %w", result.Str, err)
		}

		rows[index] = []string{appDetails.Name, appDetails.Bucket, appDetails.Version}
	}

	return output.PrintTable(headers, rows)
}
