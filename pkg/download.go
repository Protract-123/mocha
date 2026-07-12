package pkg

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	RealFilename string
}

func DownloadPackageFiles(manifestRef manifest.Ref, downloadArch string, force bool, mochaDir string) ([]DownloadResult, error) {
	downloadEntries, err := manifest.GetManifestDownloads(manifestRef.ManifestPath, downloadArch)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest downloads: %w", err)
	}

	downloadResults := make([]DownloadResult, 0, len(downloadEntries))

	for _, entry := range downloadEntries {
		downloadPath, err := fileops.GetCachePath(mochaDir, manifestRef.Name, manifestRef.Version, entry.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to get cache path: %w", err)
		}
		filename := filepath.Base(downloadPath)

		if _, err := os.Stat(downloadPath); err != nil || force {
			output.LogOutput(fmt.Sprintf("Downloading %s to %s", entry.URL, downloadPath))
			if err := fileops.DownloadFile(entry.URL, downloadPath); err != nil {
				return nil, fmt.Errorf("failed to download %s: %w", filename, err)
			}
			output.LogOutput(fmt.Sprintf("Downloaded %s", filename))
		} else {
			output.LogOutput(fmt.Sprintf("Cache hit, skipping %s", filename))
		}

		if err := fileops.VerifyHash(downloadPath, entry.Hash); err != nil {
			_ = os.Remove(downloadPath)
			return nil, fmt.Errorf("failed to verify %s: %w", filename, err)
		}

		output.LogOutput(fmt.Sprintf("Verified %s\n", filename))

		downloadResults = append(downloadResults, DownloadResult{
			Entry:        entry,
			DownloadPath: downloadPath,
			Filename:     filename,
			RealFilename: getFileNameFromUrl(entry.URL),
		})
	}

	return downloadResults, nil
}

func GetDownloadArch() (string, error) {
	cpuArch := runtime.GOARCH

	if cpuArch == "386" {
		return "32bit", nil
	} else if cpuArch == "amd64" {
		return "64bit", nil
	} else if cpuArch == "arm64" {
		return "arm64", nil
	}

	return "", fmt.Errorf("cpu architecture %q is unsupported", cpuArch)
}

func getFileNameFromUrl(rawUrl string) string {
	realFilename := path.Base(rawUrl)
	if parsedURL, parseErr := url.Parse(rawUrl); parseErr == nil {
		if strings.HasPrefix(parsedURL.Fragment, "/") {
			realFilename = path.Base(parsedURL.Fragment)
		} else {
			realFilename = path.Base(parsedURL.Path)
		}
	}

	return realFilename
}
