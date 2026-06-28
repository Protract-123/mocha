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
	query := strings.ToLower(cmd.Query)
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

		manifests, err := os.ReadDir(filepath.Join(bucketsDir, bucket.Name(), "bucket"))
		if err != nil {
			return fmt.Errorf("failed to read %s's manifest directory: %w", bucket.Name(), err)
		}

		for _, entry := range manifests {
			manifestName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

			if manifestName == query {
				exactMatches = append(exactMatches, manifestName)
				continue
			}

			allManifestNames = append(allManifestNames, manifestName)
		}
	}

	fuzzyMatches := fuzzy.Find(query, allManifestNames)

	if len(fuzzyMatches) == 0 && len(exactMatches) == 0 {
		return fmt.Errorf("no results found in buckets")
	}

	if len(exactMatches) != 0 {
		if err := outputResults(exactMatches, mochaDir); err != nil {
			return fmt.Errorf("failed to print exact matches: %w", err)
		}
	}

	if len(fuzzyMatches) != 0 {
		limit := len(fuzzyMatches)
		if cmd.Count > 0 && cmd.Count < limit {
			limit = cmd.Count
		}

		fuzzyResults := make([]string, limit)
		for index, result := range fuzzyMatches[:limit] {
			fuzzyResults[index] = result.Str
		}

		if err := outputResults(fuzzyResults, mochaDir); err != nil {
			return fmt.Errorf("failed to print fuzzy matches: %w", err)
		}
	}

	return nil
}

func outputResults(matches []string, mochaDir string) error {
	headers := []string{"Name", "Bucket", "Version"}
	rows := make([][]string, len(matches)+1)

	for index, result := range matches {
		manifestRef, err := manifest.ParseRefString(result)
		if err != nil {
			return fmt.Errorf("failed to parse result %q: %w", result, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", result, err)
		}

		rows[index] = []string{manifestRef.Name, manifestRef.Bucket, manifestRef.Version}
	}

	if err := output.PrintTable(headers, rows); err != nil {
		return fmt.Errorf("failed to output table: %w", err)
	}

	return nil
}
