package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
)

type DownloadResult struct {
	Entry        DownloadEntry
	DownloadPath string
	Filename     string
}

func DownloadManifestFiles(refString string, force bool, mochaDir string) (*Ref, []DownloadResult, error) {
	manifestRef, err := ParseRefString(refString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
	}

	manifestRef, err = PopulateRef(manifestRef, mochaDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get %q manifest details: %w", refString, err)
	}

	downloadArch, err := GetDownloadArch()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get system architecture: %w", err)
	}

	downloadEntries, err := GetManifestDownloads(manifestRef.ManifestPath, downloadArch)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get manifest downloads: %w", err)
	}

	downloadResults := make([]DownloadResult, 0, len(downloadEntries))

	for _, entry := range downloadEntries {
		downloadPath := fileops.GetCachePath(mochaDir, manifestRef.Name, manifestRef.Version, entry.URL)
		filename := filepath.Base(downloadPath)

		if _, err := os.Stat(downloadPath); err != nil || force {
			output.LogOutput(fmt.Sprintf("Downloading %s to %s", entry.URL, downloadPath))
			if err := fileops.DownloadFile(entry.URL, downloadPath); err != nil {
				return nil, nil, fmt.Errorf("failed to download %s: %w", filename, err)
			}
			output.LogOutput(fmt.Sprintf("Downloaded %s", filename))
		} else {
			output.LogOutput(fmt.Sprintf("Cache hit, skipping %s", filename))
		}

		if err := fileops.VerifyHash(downloadPath, entry.Hash); err != nil {
			_ = os.Remove(downloadPath)
			return nil, nil, fmt.Errorf("failed to verify %s: %w", filename, err)
		}

		output.LogOutput(fmt.Sprintf("Verified %s\n", filename))

		downloadResults = append(downloadResults, DownloadResult{
			Entry:        entry,
			DownloadPath: downloadPath,
			Filename:     filename,
		})
	}

	return &manifestRef, downloadResults, nil
}

func InstallManifestFile(filePath string, installDir string, subDir string, mochaDir string) error {
	extension := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filepath.Base(filePath), extension)

	tempDir := filepath.Join(mochaDir, "temp", fileName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", tempDir, err)
	}
	defer os.RemoveAll(filepath.Join(mochaDir, "temp"))

	var extractionError error

	switch extension {
	case ".zip":
		extractionError = fileops.ExtractZip(filePath, tempDir)
	case ".msi":
		extractionError = fileops.ExtractMsi(filePath, tempDir)
	default:
		// 7zip extract by default
		return nil
	}

	if extractionError != nil {
		return fmt.Errorf("failed to extract %s to %s: %w", filePath, tempDir, extractionError)
	}

	extractedDir := filepath.Join(tempDir, subDir)

	if err := mergeDir(extractedDir, installDir); err != nil {
		return fmt.Errorf("failed to merge %s into %s: %w", subDir, installDir, err)
	}

	return nil
}

func mergeDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		dstPath := filepath.Join(dst, relativePath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		if err := os.Rename(path, dstPath); err != nil {
			return fmt.Errorf("failed to move %s to %s: %w", path, dstPath, err)
		}
		return nil
	})
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
