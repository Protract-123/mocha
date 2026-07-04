package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

type CacheCommand struct {
	Show  *showCacheCommand  `arg:"subcommand:show" help:"show all cache items"`
	Clear *clearCacheCommand `arg:"subcommand:clear" help:"clear cache items"`
}

type showCacheCommand struct{}
type clearCacheCommand struct {
	Items []string `arg:"positional" help:"cache items to remove (e.g. git, bat@0.26.1)"`
	All   bool     `arg:"-a,--all" help:"clear all cache items"`
}

func (cmd *CacheCommand) Run(mochaDir string) error {
	switch {
	case cmd.Show != nil:
		return cmd.Show.Run(mochaDir)
	case cmd.Clear != nil:
		return cmd.Clear.Run(mochaDir)
	default:
		return arg.ErrHelp
	}
}

func (cmd *showCacheCommand) Run(mochaDir string) error {
	cacheItems, err := fileops.GetCacheItems(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	if len(cacheItems) == 0 {
		output.LogOutput("no cache items to show")
		return nil
	}

	type cacheItemKey struct {
		name    string
		version string
	}

	var cacheGroupOrder []cacheItemKey
	groupedCacheItems := make(map[cacheItemKey]*fileops.CacheItem)

	for _, cacheItem := range cacheItems {
		key := cacheItemKey{cacheItem.Name, cacheItem.Version}

		if existing, ok := groupedCacheItems[key]; ok {
			existing.Size += cacheItem.Size
			continue
		}

		groupedCacheItems[key] = &cacheItem
		cacheGroupOrder = append(cacheGroupOrder, key)
	}

	headers := []string{"Name", "Version", "Size"}
	rows := make([][]string, len(groupedCacheItems))

	var totalBytes int64

	for i, key := range cacheGroupOrder {
		item := groupedCacheItems[key]
		rows[i] = []string{
			item.Name,
			item.Version,
			convertToHumanReadable(item.Size),
		}

		totalBytes += item.Size
	}

	if err := output.PrintTable(headers, rows); err != nil {
		return fmt.Errorf("failed to display cache items: %w", err)
	}

	fmt.Printf("\nTotal Size: %s\n", convertToHumanReadable(totalBytes))

	return nil
}

func convertToHumanReadable(bytes int64) string {
	var units = [...]string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

	value := float64(bytes)
	unit := "Bytes"

	i := 0
	for value >= 1024 && i < len(units) {
		value = value / 1024
		unit = units[i]
		i++
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

func (cmd *clearCacheCommand) Run(mochaDir string) error {
	if len(cmd.Items) == 0 && !cmd.All {
		return arg.ErrHelp
	}

	cacheItems, err := fileops.GetCacheItems(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	if len(cmd.Items) == 0 && cmd.All {
		for _, cacheItem := range cacheItems {
			if err := os.Remove(cacheItem.Path); err != nil {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}

		return nil
	}

	for _, refString := range cmd.Items {
		cacheRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse cache item %q: %w", refString, err)
		}

		for _, cacheItem := range cacheItems {
			if cacheRef.Name != cacheItem.Name && cacheRef.Name != "" {
				continue
			}

			if cacheRef.Version != cacheItem.Version && cacheRef.Version != "" {
				continue
			}

			if err := os.Remove(cacheItem.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}
	}

	return nil
}
