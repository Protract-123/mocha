package show

import (
	"fmt"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
)

type Command struct{}

func (cmd *Command) Run() error {
	cacheItems, err := fileops.GetCacheItems(config.Current().MochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	if len(cacheItems) == 0 {
		return fmt.Errorf("no cache items to show")
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

	tableConfig := output.TableConfig{
		Spacing: 2,
		Alignments: []output.Alignment{
			output.LeftAlign,
			output.LeftAlign,
			output.RightAlign,
		},
		BorderStyle: output.LightBorder,
	}

	if err := output.PrintTable(headers, rows, tableConfig); err != nil {
		return fmt.Errorf("failed to display cache items: %w", err)
	}

	output.LogOutput("Total Size: %s", convertToHumanReadable(totalBytes))

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
