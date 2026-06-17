package bucket

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Metadata struct {
	Name          string
	Source        string
	LastUpdated   string
	ManifestCount int
}

func GetAllBucketMetadata(mochaDir string) ([]Metadata, error) {
	bucketsDir := filepath.Join(mochaDir, "buckets")

	buckets, err := os.ReadDir(bucketsDir)
	if err != nil {
		return nil, err
	}

	var bucketMetadata []Metadata

	for _, entry := range buckets {
		if !entry.IsDir() {
			continue
		}

		metadata, err := GetBucketMetadata(mochaDir, entry.Name())
		if err != nil {
			return nil, err
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
		return Metadata{}, err
	}
	bucketSource := strings.TrimSpace(string(sourceOut))

	updatedCmd := exec.Command("git", "log", "--format=%aD", "-n", "1")
	updatedCmd.Dir = bucketPath
	updatedOut, err := updatedCmd.Output()
	if err != nil {
		return Metadata{}, err
	}

	bucketLastUpdated, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimSpace(string(updatedOut)))
	if err != nil {
		return Metadata{}, err
	}

	manifests, err := os.ReadDir(filepath.Join(bucketPath, "bucket"))
	if err != nil {
		return Metadata{}, err
	}
	manifestCount := len(manifests)

	return Metadata{
		Name:          bucketName,
		Source:        bucketSource,
		LastUpdated:   bucketLastUpdated.Format("02-01-2006 15:04:05"),
		ManifestCount: manifestCount,
	}, nil
}
