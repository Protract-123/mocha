package bucket

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type knownCmd struct{}

const scoopBucketsJsonSource = "https://raw.githubusercontent.com/ScoopInstaller/Scoop/refs/heads/master/buckets.json"

func (cmd *knownCmd) Run() error {
	_, err := os.Stat("buckets.json")
	if os.IsNotExist(err) {
		err := updateKnownBuckets()
		if err != nil {
			return err
		}
	}

	bucketsJson, err := os.ReadFile("buckets.json")
	if err != nil {
		return err
	}

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

		fmt.Print("\033[1;35m", name, "\033[0m: ", url, "\n")
	}

	return nil
}

// TODO: move to update command

func updateKnownBuckets() error {
	exePath, err := os.Executable()
	if err != nil {
		println(err)
	}

	exeDir := filepath.Dir(exePath)
	bucketsPath := filepath.Join(exeDir, "buckets.json")

	err = downloadFile(scoopBucketsJsonSource, bucketsPath)
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
