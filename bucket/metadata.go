package bucket

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Metadata struct {
	Name          string
	Source        string
	LastUpdated   time.Time
	ManifestCount int
}

func GetAllBucketMetadata(mochaDir string) ([]Metadata, error) {
	buckets, err := os.ReadDir(filepath.Join(mochaDir, "buckets"))
	if err != nil {
		return nil, fmt.Errorf("failed to get all buckets: %w", err)
	}

	bucketMetadata := make([]Metadata, 0, len(buckets))
	for _, bucket := range buckets {
		if !bucket.IsDir() {
			continue
		}

		metadata, err := GetBucketMetadata(mochaDir, bucket.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for bucket %q: %w", bucket.Name(), err)
		}

		bucketMetadata = append(bucketMetadata, metadata)
	}

	return bucketMetadata, nil
}

func GetBucketMetadata(mochaDir string, bucketName string) (Metadata, error) {
	bucketPath := filepath.Join(mochaDir, "buckets", bucketName)

	sourceCmd := exec.Command("git", "config", "remote.origin.url")
	sourceCmd.Dir = bucketPath
	sourceOut, err := sourceCmd.Output()
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to get bucket source: %w", err)
	}
	bucketSource := strings.TrimSpace(string(sourceOut))

	updatedCmd := exec.Command("git", "log", "--format=%aD", "-n", "1")
	updatedCmd.Dir = bucketPath
	updatedOut, err := updatedCmd.Output()
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to get last update date: %w", err)
	}

	bucketLastUpdated, err := time.Parse(time.RFC1123Z, strings.TrimSpace(string(updatedOut)))
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to parse last update date: %w", err)
	}

	manifests, err := os.ReadDir(filepath.Join(bucketPath, "bucket"))
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to get manifest count: %w", err)
	}

	return Metadata{
		Name:          bucketName,
		Source:        bucketSource,
		LastUpdated:   bucketLastUpdated,
		ManifestCount: len(manifests),
	}, nil
}
