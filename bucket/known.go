package bucket

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type knownCmd struct{}

const scoopBucketsJsonSource = "https://raw.githubusercontent.com/ScoopInstaller/Scoop/refs/heads/master/buckets.json"

func (cmd *knownCmd) Run(mochaDir string) error {
	knownBucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	if !filepath.IsAbs(knownBucketsPath) {
		return fmt.Errorf("%s is not an absolute path", knownBucketsPath)
	}

	_, err := os.Stat(knownBucketsPath)
	if os.IsNotExist(err) {
		err := updateKnownBuckets(mochaDir)
		if err != nil {
			return err
		}
	}

	knownBuckets, err := parseBucketList(knownBucketsPath)
	if err != nil {
		return err
	}

	for _, bucket := range knownBuckets {
		fmt.Print("\033[1;35m", bucket.Name, "\033[0m: ", bucket.URL, "\n")
	}

	return nil
}

// TODO: move to update command

func updateKnownBuckets(mochaDir string) error {
	bucketsPath := filepath.Join(mochaDir, "known_buckets.json")

	err := downloadFile(scoopBucketsJsonSource, bucketsPath)
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
