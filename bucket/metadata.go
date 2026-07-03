package bucket

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type Metadata struct {
	Name          string
	Source        string
	LastUpdated   time.Time
	ManifestCount int
}

func GetAllBucketMetadata(mochaDir string) ([]Metadata, error) {
	dirs, err := os.ReadDir(filepath.Join(mochaDir, "buckets"))
	if err != nil {
		return nil, fmt.Errorf("failed to get all buckets: %w", err)
	}

	var buckets []os.DirEntry
	for _, bucket := range dirs {
		if bucket.IsDir() {
			buckets = append(buckets, bucket)
		}
	}

	group, ctx := errgroup.WithContext(context.Background())
	bucketMetadata := make([]Metadata, len(buckets))

	for index, bucket := range buckets {
		group.Go(func() error {
			metadata, err := GetBucketMetadata(mochaDir, bucket.Name(), ctx)
			if err != nil {
				return fmt.Errorf("failed to get metadata for bucket %q: %w", bucket.Name(), err)
			}

			bucketMetadata[index] = metadata
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return bucketMetadata, nil
}

func GetBucketMetadata(mochaDir string, bucketName string, ctx context.Context) (Metadata, error) {
	bucketPath := filepath.Join(mochaDir, "buckets", bucketName)

	sourceCmd := exec.CommandContext(ctx, "git", "config", "remote.origin.url")
	sourceCmd.Dir = bucketPath
	sourceOut, err := sourceCmd.Output()
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to get bucket source: %w", err)
	}
	bucketSource := strings.TrimSpace(string(sourceOut))

	updatedCmd := exec.CommandContext(ctx, "git", "log", "--format=%aI", "-n", "1")
	updatedCmd.Dir = bucketPath
	updatedOut, err := updatedCmd.Output()
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to get last update date: %w", err)
	}

	bucketLastUpdated, err := time.Parse(time.RFC3339, strings.TrimSpace(string(updatedOut)))
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
