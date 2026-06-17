package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DownloadEntry struct {
	URL  string
	Hash string
}

func GetManifestPath(bucketName string, manifestName string, mochaDir string) (string, error) {
	manifestPath := filepath.Join(mochaDir, "buckets", bucketName, "bucket", fmt.Sprintf("%s.json", manifestName))

	_, err := os.Stat(manifestPath)
	if err != nil {
		return "", err
	}

	return manifestPath, nil
}

func GetManifestDownloads(manifestPath string, architecture string) ([]DownloadEntry, error) {
	rawData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]any
	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return nil, err
	}

	var urls, hashes []string

	if archMap, ok := jsonData["architecture"].(map[string]any); ok {
		if archBlock, ok := archMap[architecture].(map[string]any); ok {
			urls = extractStringOrArray(archBlock["url"])
			hashes = extractStringOrArray(archBlock["hash"])
		}
	}

	if len(urls) == 0 {
		urls = extractStringOrArray(jsonData["url"])
		if len(hashes) == 0 {
			hashes = extractStringOrArray(jsonData["hash"])
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("unable to find download URL for %s (arch %q)", manifestPath, architecture)
	}

	entries := make([]DownloadEntry, len(urls))
	for i, u := range urls {
		entry := DownloadEntry{URL: u}
		if i < len(hashes) {
			entry.Hash = hashes[i]
		}
		entries[i] = entry
	}

	return entries, nil
}

func GetManifestVersion(manifestPath string) (string, error) {
	rawData, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	var jsonData map[string]any
	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return "", err
	}

	if version, ok := jsonData["version"].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("unable to find app version for %s", manifestPath)
}

func extractStringOrArray(v any) []string {
	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		return []string{val}
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}
