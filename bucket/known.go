package bucket

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const knownBucketsSourceFile = "https://raw.githubusercontent.com/ScoopInstaller/Scoop/refs/heads/master/buckets.json"

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

	err := downloadFile(knownBucketsSourceFile, bucketsPath)
	if err != nil {
		return err
	}

	return nil
}

// Taken from https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			println(err)
		}
	}(out)

	client := &http.Client{Timeout: 30 * time.Second}

	// Get the data
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
