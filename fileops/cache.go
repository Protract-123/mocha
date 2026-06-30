package fileops

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type CacheItem struct {
	Name    string
	Version string
	Size    int64
	Path    string
}

func GetCachePath(mochaDir string, app string, version string, rawURL string) (string, error) {
	sum := sha256.Sum256([]byte(rawURL))
	shortHash := hex.EncodeToString(sum[:])[:7]
	var ext string

	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		ext = filepath.Ext(rawURL)
	} else if strings.HasPrefix(parsedUrl.Fragment, "/") {
		ext = filepath.Ext(parsedUrl.Fragment)
	} else {
		ext = filepath.Ext(parsedUrl.Path)
	}

	if ext == "" {
		return "", fmt.Errorf("failed to determine file extension for %s", rawURL)
	}

	return filepath.Join(mochaDir, "cache", fmt.Sprintf("%s#%s#%s%s", app, version, shortHash, ext)), nil
}

func GetCacheItems(mochaDir string) ([]CacheItem, error) {
	cacheDir := filepath.Join(mochaDir, "cache")

	cacheFiles, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	cacheItems := make([]CacheItem, 0, len(cacheFiles))

	for _, item := range cacheFiles {
		if item.IsDir() {
			continue
		}

		fileParts := strings.Split(item.Name(), "#")
		if len(fileParts) != 3 {
			continue
		}

		fileInfo, err := item.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to read %s info: %w", item.Name(), err)
		}

		name := fileParts[0]
		version := fileParts[1]
		fileSize := fileInfo.Size()

		cacheItems = append(cacheItems, CacheItem{name, version, fileSize, filepath.Join(cacheDir, item.Name())})
	}

	return cacheItems, nil
}
