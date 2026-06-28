package commands

import (
	"fmt"

	"github.com/Protract-123/mocha/bucket"
)

type UpdateCommand struct {
	Buckets []string `arg:"positional"`
}

func (cmd *UpdateCommand) Run(mochaDir string) error {
	if err := bucket.UpdateKnownBuckets(mochaDir); err != nil {
		return fmt.Errorf("failed to update known buckets: %w", err)
	}

	if len(cmd.Buckets) == 0 {
		if err := bucket.UpdateAllBuckets(mochaDir); err != nil {
			return fmt.Errorf("failed to update all buckets: %w", err)
		}

		return nil
	}

	for _, entry := range cmd.Buckets {
		if err := bucket.UpdateBucket(entry, mochaDir); err != nil {
			return fmt.Errorf("failed to update bucket %q: %w", entry, err)
		}
	}

	return nil
}
