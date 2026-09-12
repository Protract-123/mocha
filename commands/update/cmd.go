package update

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
)

type Command struct {
	Buckets []string `arg:"positional" help:"buckets to update; updates all buckets if omitted (e.g. main)"`
}

func (cmd *Command) Run() error {
	mochaDir := config.Current().MochaDirectory

	if err := bucket.UpdateKnownBuckets(mochaDir); err != nil {
		return fmt.Errorf("failed to update known buckets: %w", err)
	}

	var buckets []string
	if len(cmd.Buckets) == 0 {
		entries, err := os.ReadDir(filepath.Join(mochaDir, "buckets"))
		if err != nil {
			return fmt.Errorf("failed to get all buckets: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				buckets = append(buckets, entry.Name())
			}
		}
	} else {
		buckets = cmd.Buckets
	}

	group := sync.WaitGroup{}

	for _, entry := range buckets {
		group.Go(func() {
			if err := bucket.Update(entry, mochaDir); err != nil {
				output.LogError(fmt.Errorf("failed to update bucket %q: %w", entry, err))
			} else {
				output.LogSuccess("successfully updated bucket %q", entry)
			}
		})
	}

	group.Wait()
	return nil
}
