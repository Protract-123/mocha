package bucket

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Bucket struct {
	Name   string `toml:"name"`
	Source string `toml:"source"`
}

func DownloadBucket(bucket Bucket, mochaDir string) error {
	bucketsDir := filepath.Join(mochaDir, "buckets")

	if err := os.MkdirAll(bucketsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to make bucket directory: %w", err)
	}

	destDir := filepath.Join(bucketsDir, bucket.Name)

	_, err := os.Stat(destDir)
	if err == nil {
		return fmt.Errorf("bucket %s already exists", bucket.Name)
	}

	cmd := exec.Command("git", "clone", bucket.Source, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to clone bucket: %w", err)
	}
	return nil
}

func DeleteBucket(name string, mochaDir string) error {
	bucketDir := filepath.Join(mochaDir, "buckets", name)

	err := os.RemoveAll(bucketDir)
	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %v", name, err)
	}
	return nil
}

func ParseBucketList(file string) ([]Bucket, error) {
	bucketsJson, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read buckets list file: %q", err)
	}

	var buckets []Bucket

	for _, line := range strings.Split(string(bucketsJson), "\n") {
		line = strings.TrimSpace(line)

		if line == "{" || line == "}" || line == "" {
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.Trim(parts[0], `"`)
		url := strings.Trim(parts[1], `",`)

		buckets = append(buckets, Bucket{Name: name, Source: url})
	}

	return buckets, nil
}
