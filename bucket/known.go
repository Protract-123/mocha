package bucket

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/fileops"
)

const knownBucketsSourceFile = "https://raw.githubusercontent.com/ScoopInstaller/Scoop/refs/heads/master/buckets.json"

func GetKnownBucket(name string, mochaDir string) (Bucket, error) {
	knownBuckets, err := GetKnownBuckets(mochaDir)
	if err != nil {
		return Bucket{}, fmt.Errorf("failed to get known buckets: %w", err)
	}

	for _, entry := range knownBuckets {
		if entry.Name == name {
			return entry, nil
		}
	}

	return Bucket{}, fmt.Errorf("bucket %q is not a known bucket", name)
}

func GetKnownBuckets(mochaDir string) ([]Bucket, error) {
	knownBucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	if _, err := os.Stat(knownBucketsPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := UpdateKnownBuckets(mochaDir); err != nil {
				return nil, fmt.Errorf("failed to update known buckets: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to check if known_buckets.json exists: %w", err)
		}
	}

	knownBuckets, err := ParseBucketList(knownBucketsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse known buckets: %w", err)
	}

	return knownBuckets, nil
}

func UpdateKnownBuckets(mochaDir string) error {
	knownBucketsPath := filepath.Join(mochaDir, "known_buckets.json")
	if err := fileops.DownloadFile(knownBucketsSourceFile, knownBucketsPath); err != nil {
		return fmt.Errorf("failed to download known_buckets.json: %w", err)
	}

	return nil
}
