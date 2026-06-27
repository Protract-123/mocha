package commands

import (
	"fmt"
	"os"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

type CacheCommand struct {
	List  *listCacheCommand  `arg:"subcommand:list"`
	Clear *clearCacheCommand `arg:"subcommand:clear"`
}

type listCacheCommand struct{}
type clearCacheCommand struct {
	ManifestReferences []string `arg:"positional"`
}

func (cmd CacheCommand) Run(mochaDir string) error {
	if cmd.List != nil {
		err := cmd.List.Run(mochaDir)
		if err != nil {
			return fmt.Errorf("failed to list cache items: %w", err)
		}
	}
	if cmd.Clear != nil {
		err := cmd.Clear.Run(mochaDir)
		if err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	return nil
}

func (cmd listCacheCommand) Run(mochaDir string) error {
	rawCacheItems, err := fileops.GetCacheItems(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	type cacheItemKey struct{ name, version string }
	cacheItems := make(map[cacheItemKey]*fileops.CacheItem)
	var cacheItemOrder []cacheItemKey

	for _, item := range rawCacheItems {
		key := cacheItemKey{item.Name, item.Version}
		if existing, ok := cacheItems[key]; ok {
			existing.Size += item.Size
		} else {
			cacheItems[key] = &item
			cacheItemOrder = append(cacheItemOrder, key)
		}
	}

	headers := []string{"Name", "Version", "Size"}
	rows := make([][]string, len(cacheItems))

	var totalBytes int64 = 0

	for i, key := range cacheItemOrder {
		item := cacheItems[key]
		rows[i] = []string{
			item.Name,
			item.Version,
			ConvertToHumanReadable(item.Size),
		}

		totalBytes += item.Size
	}

	err = output.PrintTable(headers, rows)
	if err != nil {
		return fmt.Errorf("failed to display cache items: %w", err)
	}

	fmt.Printf("\nTotal Size: %s\n", ConvertToHumanReadable(totalBytes))

	return nil
}

func ConvertToHumanReadable(bytes int64) string {
	var units = [...]string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

	value := float32(bytes)
	unit := "Bytes"

	i := 0
	for value >= 1024 {
		value = value / 1024
		unit = units[i]
		i += 1
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

func (cmd clearCacheCommand) Run(mochaDir string) error {
	cacheItems, err := fileops.GetCacheItems(mochaDir)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	if len(cmd.ManifestReferences) == 0 {
		for _, cacheItem := range cacheItems {
			err := os.Remove(cacheItem.Path)
			if err != nil {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}

		return nil
	}

	for _, refString := range cmd.ManifestReferences {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		for _, cacheItem := range cacheItems {
			if manifestRef.Name != cacheItem.Name && manifestRef.Name != "" {
				continue
			}

			if manifestRef.Version != cacheItem.Version && manifestRef.Version != "" {
				continue
			}

			err := os.Remove(cacheItem.Path)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}
	}

	return nil
}
