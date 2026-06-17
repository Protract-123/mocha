package fileops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Taken from https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9

func DownloadFile(url string, downloadPath string) error {
	err := os.MkdirAll(filepath.Dir(downloadPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(downloadPath)
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
