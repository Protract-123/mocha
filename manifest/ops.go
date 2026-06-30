package manifest

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
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
	RealFilename string
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
		downloadPath, err := fileops.GetCachePath(mochaDir, manifestRef.Name, manifestRef.Version, entry.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get cache path: %w", err)
		}
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

		realFilename := path.Base(entry.URL)
		if parsedURL, parseErr := url.Parse(entry.URL); parseErr == nil {
			if strings.HasPrefix(parsedURL.Fragment, "/") {
				realFilename = path.Base(parsedURL.Fragment)
			} else {
				realFilename = path.Base(parsedURL.Path)
			}
		}

		downloadResults = append(downloadResults, DownloadResult{
			Entry:        entry,
			DownloadPath: downloadPath,
			Filename:     filename,
			RealFilename: realFilename,
		})
	}

	return &manifestRef, downloadResults, nil
}

type InstallOptions struct {
	SubDir string

	InnoSetup    bool
	RealFileName string
}

func InstallManifestFile(filePath string, installDir string, mochaDir string, options InstallOptions) error {
	extension := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filepath.Base(filePath), extension)

	tempDir := filepath.Join(mochaDir, "temp", fileName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", tempDir, err)
	}
	defer os.RemoveAll(filepath.Join(mochaDir, "temp"))

	switch extension {
	case ".zip":
		if err := fileops.ExtractZip(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	case ".msi":
		if err := fileops.ExtractMsi(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	case ".exe":
		if !options.InnoSetup {
			if options.RealFileName == "" {
				return fmt.Errorf("missing real file name for %s", filePath)
			}
			targetFilePath := filepath.Join(tempDir, options.RealFileName)

			if err := copyFile(filePath, targetFilePath); err != nil {
				return fmt.Errorf("failed to move %s: %w", targetFilePath, err)
			}
		} else {
			if err := fileops.ExtractInnoSetup(filePath, tempDir, options.SubDir); err != nil {
				return fmt.Errorf("failed to extract InnoSetup from %s: %w", filePath, err)
			}
			options.SubDir = ""
		}
	default:
		if err := fileops.Extract7z(filePath, tempDir); err != nil {
			return fmt.Errorf("failed to extract %s: %w", filePath, err)
		}
	}

	extractedDir := filepath.Join(tempDir, options.SubDir)
	if err := mergeDir(extractedDir, installDir); err != nil {
		return fmt.Errorf("failed to merge %s into %s: %w", extractedDir, installDir, err)
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

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", src, err)
	}
	defer srcFile.Close()

	targetFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", dst, err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file %s: %w", dst, err)
	}

	if err := targetFile.Sync(); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}

	return nil
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
