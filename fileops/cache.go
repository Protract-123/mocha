package fileops

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
)

func GetCachePath(mochaDir string, app string, version string, url string) string {
	sum := sha256.Sum256([]byte(url))
	shortHash := hex.EncodeToString(sum[:])[:7]
	ext := extractFileExtension(url)

	return filepath.Join(mochaDir, "cache", fmt.Sprintf("%s#%s#%s%s", app, version, shortHash, ext))
}

func GetCacheItem() {

}

func extractFileExtension(rawURL string) string {
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		return filepath.Ext(rawURL)
	}
	return filepath.Ext(parsedUrl.Path)
}
