package bucket

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/fileops"
)

const knownBucketsSourceFile = "https://raw.githubusercontent.com/ScoopInstaller/Scoop/refs/heads/master/buckets.json"

func GetKnownBucket(name string, mochaDir string) (Bucket, error) {
	knownBuckets, err := GetKnownBuckets(mochaDir)
	if err != nil {
		return Bucket{}, fmt.Errorf("failed to get known buckets: %v", err)
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

	_, err := os.Stat(knownBucketsPath)
	if os.IsNotExist(err) {
		err := UpdateKnownBuckets(mochaDir)
		if err != nil {
			return nil, fmt.Errorf("failed to update known buckets: %w", err)
		}
	}

	knownBuckets, err := ParseBucketList(knownBucketsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse known buckets: %w", err)
	}

	return knownBuckets, nil
}

func UpdateKnownBuckets(mochaDir string) error {
	bucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	err := fileops.DownloadFile(knownBucketsSourceFile, bucketsPath)
	if err != nil {
		return fmt.Errorf("failed to download known buckets: %w", err)
	}

	return nil
}
