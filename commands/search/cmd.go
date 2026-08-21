package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/sahilm/fuzzy"
)

type Command struct {
	Query string `arg:"positional,required" help:"app name to search for"`
	Count int    `default:"20" arg:"-c, --count" help:"max number of fuzzy results to show"`
}

func (cmd *Command) Run() error {
	query := strings.ToLower(cmd.Query)

	mochaDir := config.Current().MochaDirectory
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
			return fmt.Errorf("failed to read manifest directory of %s: %w", bucket.Name(), err)
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

	limit := len(fuzzyMatches) + len(exactMatches)
	if cmd.Count > 0 && cmd.Count < limit {
		limit = cmd.Count
	}

	allMatches := make([]string, limit)

	for i := 0; i < limit && i < len(exactMatches); i++ {
		allMatches[i] = exactMatches[i]
	}

	for i := len(exactMatches); i < limit; i++ {
		allMatches[i] = fuzzyMatches[i-len(exactMatches)].Str
	}

	if err := outputResults(allMatches, mochaDir); err != nil {
		return fmt.Errorf("failed to print search results: %w", err)
	}

	return nil
}

func outputResults(matches []string, mochaDir string) error {
	headers := []string{"Name", "Bucket", "Version"}
	rows := make([][]string, len(matches))

	for index, result := range matches {
		info, err := manifest.ParseSpec(result)
		if err != nil {
			return fmt.Errorf("failed to parse result %q: %w", result, err)
		}

		info, err = manifest.PopulateInfo(info, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get manifest details for %q: %w", result, err)
		}

		rows[index] = []string{info.Name, info.Bucket, info.Version}
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
		return fmt.Errorf("failed to output table: %w", err)
	}

	return nil
}
