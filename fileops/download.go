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

	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status %q", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync downloaded file: %w", err)
	}

	return nil
}
