package clear

import (
	"errors"
	"fmt"
	"os"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

type Command struct {
	Items []string `arg:"positional" help:"cache items to remove (e.g. git, bat@0.26.1)"`
	All   bool     `arg:"-a,--all" help:"clear all cache items"`
}

func (cmd *Command) Run() error {
	if len(cmd.Items) == 0 && !cmd.All {
		return arg.ErrHelp
	}

	cacheItems, err := fileops.GetCacheItems(config.Current().MochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to get cache items: %w", err)
	}

	if len(cmd.Items) == 0 && cmd.All {
		for _, cacheItem := range cacheItems {
			if err := os.Remove(cacheItem.Path); err != nil {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}

		output.LogSuccess("successfully deleted all cache items")
		return nil
	}

	for _, spec := range cmd.Items {
		cacheInfo, err := manifest.ParseSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse cache item %q: %w", spec, err)
		}

		for _, cacheItem := range cacheItems {
			if cacheInfo.Name != cacheItem.Name && cacheInfo.Name != "" {
				continue
			}

			if cacheInfo.Version != cacheItem.Version && cacheInfo.Version != "" {
				continue
			}

			if err := os.Remove(cacheItem.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to remove cache item %q: %w", cacheItem.Path, err)
			}
		}

		output.LogSuccess("successfully deleted cache for %q", spec)
	}

	return nil
}
