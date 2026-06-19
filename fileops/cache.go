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

func GetCachePath(mochaDir string, app string, version string, url string) string {
	sum := sha256.Sum256([]byte(url))
	shortHash := hex.EncodeToString(sum[:])[:7]
	ext := extractFileExtension(url)

	return filepath.Join(mochaDir, "cache", fmt.Sprintf("%s#%s#%s%s", app, version, shortHash, ext))
}

func GetCacheItems(mochaDir string) ([]CacheItem, error) {
	cacheDir := filepath.Join(mochaDir, "cache")

	cacheFiles, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	var cacheItems []CacheItem

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
			return nil, err
		}

		name := fileParts[0]
		version := fileParts[1]
		fileSize := fileInfo.Size()

		cacheItems = append(cacheItems, CacheItem{name, version, fileSize, filepath.Join(cacheDir, item.Name())})
	}

	return cacheItems, nil
}

func extractFileExtension(rawURL string) string {
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		return filepath.Ext(rawURL)
	}
	return filepath.Ext(parsedUrl.Path)
}
