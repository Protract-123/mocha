package bucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Bucket struct {
	Name   string
	Source string
}

func Download(bucket Bucket, mochaDir string) error {
	bucketsDir := filepath.Join(mochaDir, "buckets")
	if err := os.MkdirAll(bucketsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to make bucket directory: %w", err)
	}

	destDir := filepath.Join(bucketsDir, bucket.Name)
	if _, err := os.Stat(destDir); err == nil {
		return fmt.Errorf("bucket %s already exists", bucket.Name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check if bucket already exists: %w", err)
	}

	cmd := exec.Command("git", "clone", bucket.Source, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone bucket: %w", err)
	}

	return nil
}

func Delete(name string, mochaDir string) error {
	if err := os.RemoveAll(filepath.Join(mochaDir, "buckets", name)); err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", name, err)
	}
	return nil
}

func Update(bucketName string, mochaDir string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = filepath.Join(mochaDir, "buckets", bucketName)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	var exitErr *exec.ExitError

	switch {
	case err == nil:
		break
	case errors.As(err, &exitErr):
		cmdError := fmt.Errorf("%s: %w", strings.TrimSpace(stderr.String()), err)
		return fmt.Errorf("failed to run git pull: %w", cmdError)
	default:
		return fmt.Errorf("failed to run git pull: %w", err)

	}

	return nil
}

func ParseList(file string) ([]Bucket, error) {
	bucketsJson, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read buckets list file: %w", err)
	}

	var jsonContent map[string]string
	if err := json.Unmarshal(bucketsJson, &jsonContent); err != nil {
		return nil, fmt.Errorf("failed to parse %q: %w", file, err)
	}

	buckets := make([]Bucket, 0, len(jsonContent))
	for name, url := range jsonContent {
		buckets = append(buckets, Bucket{Name: name, Source: url})
	}

	return buckets, nil
}
