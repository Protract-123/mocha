package bucket

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type addCmd struct {
	Name          string  `arg:"positional,required"`
	RepositoryURL url.URL `arg:"positional"`
}

func (cmd *addCmd) Run(mochaDir string) error {
	knownBuckets, err := parseBucketList("buckets.json")
	if err != nil {
		return err
	}

	bucket, ok := findBucket(knownBuckets, cmd.Name)

	if !ok && cmd.RepositoryURL.String() == "" {
		return fmt.Errorf("bucket %s is not known, please provide URL", cmd.Name)
	} else if !ok && IsValidURL(cmd.RepositoryURL) {
		bucket = Bucket{
			Name: cmd.Name,
			URL:  cmd.RepositoryURL.String(),
		}
	} else if !ok {
		return fmt.Errorf("please provide valid URL for bucket %s", cmd.Name)
	}

	_, err = exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git was not found in PATH, please install git by running `mocha install git`")
	}

	bucketsDir := filepath.Join(mochaDir, "buckets")

	return getBucket(bucket, bucketsDir)
}

func findBucket(buckets []Bucket, name string) (Bucket, bool) {
	for _, bucket := range buckets {
		if bucket.Name == name {
			return bucket, true
		}
	}
	return Bucket{}, false
}

func IsValidURL(url url.URL) bool {
	if url.Scheme != "http" && url.Scheme != "https" {
		return false
	}
	if url.Host == "" {
		return false
	}

	return true
}

func getBucket(bucket Bucket, bucketsDir string) error {
	if err := os.MkdirAll(bucketsDir, os.ModePerm); err != nil {
		return err
	}

	destDir := filepath.Join(bucketsDir, bucket.Name)

	_, err := os.Stat(destDir)
	if err == nil {
		return fmt.Errorf("bucket %s already exists", bucket.Name)
	}

	cmd := exec.Command("git", "clone", bucket.URL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
