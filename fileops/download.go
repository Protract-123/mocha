package fileops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Inspired from https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9

func DownloadFile(url string, downloadPath string, progressBar bool) error {
	err := os.MkdirAll(filepath.Dir(downloadPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", downloadPath, err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ResponseHeaderTimeout = 10 * time.Second

	client := &http.Client{Transport: transport}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status %q", resp.Status)
	}

	out, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", downloadPath, err)
	}
	defer out.Close()

	fmt.Fprint(os.Stderr, "\x1b[?25l")
	defer fmt.Fprint(os.Stderr, "\x1b[?25h")

	if progressBar {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			fmt.Sprintf("downloading %s", filepath.Base(downloadPath)),
		)

		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
		if err != nil {
			_ = os.Remove(downloadPath)
			return fmt.Errorf("failed to write file: %w", err)
		}
	} else {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			_ = os.Remove(downloadPath)
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	err = out.Sync()
	if err != nil {
		_ = os.Remove(downloadPath)
		return fmt.Errorf("failed to sync downloaded file: %w", err)
	}

	return nil
}
