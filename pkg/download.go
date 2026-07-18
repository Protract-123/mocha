package pkg

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
)

type DownloadResult struct {
	Entry        manifest.DownloadEntry
	DownloadPath string
	Filename     string
}

type DownloadOptions struct {
	Force      bool
	SkipVerify bool
}

func Download(pkg Package, mochaDir string, options DownloadOptions) ([]DownloadResult, error) {
	downloadEntries, err := manifest.GetDownloadEntries(pkg.Json, pkg.Arch)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest downloads: %w", err)
	}

	downloadResults := make([]DownloadResult, 0, len(downloadEntries))

	for _, entry := range downloadEntries {
		filename := getFileNameFromUrl(entry.URL)

		downloadPath, err := fileops.GetCachePath(mochaDir, pkg.Name, pkg.Version, entry.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to get cache path: %w", err)
		}

		cached, err := isCached(downloadPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check cache path: %w", err)
		}

		if !cached || options.Force {
			output.LogOutput("Downloading %s to %s", entry.URL, downloadPath)
			if err := fileops.DownloadFile(entry.URL, downloadPath); err != nil {
				return nil, fmt.Errorf("failed to download %s: %w", filename, err)
			}
			output.LogOutput("Downloaded %s", filename)
		} else {
			output.LogOutput("Cache hit, skipping %s", filename)
		}

		if !options.SkipVerify {
			if err := fileops.VerifyHash(downloadPath, entry.Hash); err != nil {
				_ = os.Remove(downloadPath)
				return nil, fmt.Errorf("failed to verify %s: %w", filename, err)
			}
			output.LogOutput("Verified %s\n", filename)
		}

		downloadResults = append(downloadResults, DownloadResult{
			Entry:        entry,
			DownloadPath: downloadPath,
			Filename:     filename,
		})
	}

	return downloadResults, nil
}

func DownloadArch() (string, error) {
	switch runtime.GOARCH {
	case "386":
		return "32bit", nil
	case "amd64":
		return "64bit", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("cpu architecture %q is unsupported", runtime.GOARCH)
	}
}

func getFileNameFromUrl(rawUrl string) string {
	filename := path.Base(rawUrl)
	if parsedURL, parseErr := url.Parse(rawUrl); parseErr == nil {
		if strings.HasPrefix(parsedURL.Fragment, "/") {
			filename = path.Base(parsedURL.Fragment)
		} else {
			filename = path.Base(parsedURL.Path)
		}
	}

	return filename
}

func isCached(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}
