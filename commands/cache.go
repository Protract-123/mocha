package commands

import (
	"fmt"
	"os"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
)

type CacheCommand struct {
	List  *listCacheCommand  `arg:"subcommand:list"`
	Clear *clearCacheCommand `arg:"subcommand:clear"`
}

type listCacheCommand struct{}
type clearCacheCommand struct {
	Apps []string `arg:"positional"`
}

func (cmd CacheCommand) Run(mochaDir string) error {
	if cmd.List != nil {
		return cmd.List.Run(mochaDir)
	}
	if cmd.Clear != nil {
		return cmd.Clear.Run(mochaDir)
	}

	return nil
}

func (cmd listCacheCommand) Run(mochaDir string) error {
	rawCacheItems, err := fileops.GetCacheItems(mochaDir)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	if len(cmd.Apps) == 0 {
		for _, cacheItem := range cacheItems {
			err := os.Remove(cacheItem.Path)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, appString := range cmd.Apps {
		appRef, err := app.ParseAppString(appString)
		if err != nil {
			return err
		}

		for _, cacheItem := range cacheItems {
			if appRef.Name != cacheItem.Name && appRef.Name != "" {
				continue
			}

			if appRef.Version != cacheItem.Version && appRef.Version != "" {
				continue
			}

			err := os.Remove(cacheItem.Path)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}
