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
		return Bucket{}, err
	}

	for _, entry := range knownBuckets {
		if entry.Name == name {
			return entry, nil
		}
	}
	return Bucket{}, fmt.Errorf("bucket %s not found", name)
}

func GetKnownBuckets(mochaDir string) ([]Bucket, error) {
	knownBucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	if !filepath.IsAbs(knownBucketsPath) {
		return nil, fmt.Errorf("%s is not an absolute path", knownBucketsPath)
	}

	_, err := os.Stat(knownBucketsPath)
	if os.IsNotExist(err) {
		err := UpdateKnownBuckets(mochaDir)
		if err != nil {
			return nil, err
		}
	}

	knownBuckets, err := ParseBucketList(knownBucketsPath)
	if err != nil {
		return nil, err
	}

	return knownBuckets, nil
}

func UpdateKnownBuckets(mochaDir string) error {
	bucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	if err := os.MkdirAll(mochaDir, os.ModePerm); err != nil {
		return err
	}

	err := fileops.DownloadFile(knownBucketsSourceFile, bucketsPath)
	if err != nil {
		return err
	}

	return nil
}
