package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

type DownloadEntry struct {
	URL    string
	Hash   string
	SubDir string
}

func GetManifestDownloads(manifestPath string, architecture string) ([]DownloadEntry, error) {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest json: %w", err)
	}

	var urls, hashes, subDirs []string

	if archMap, ok := jsonData["architecture"].(map[string]any); ok {
		if archBlock, ok := archMap[architecture].(map[string]any); ok {
			urls = extractStringOrArray(archBlock["url"])
			hashes = extractStringOrArray(archBlock["hash"])
			subDirs = extractStringOrArray(archBlock["extract_dir"])
		}
	}

	if len(urls) == 0 {
		urls = extractStringOrArray(jsonData["url"])
		hashes = extractStringOrArray(jsonData["hash"])
		subDirs = extractStringOrArray(jsonData["extract_dir"])
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("unable to find download URL for %s (arch %q)", manifestPath, architecture)
	}

	if len(hashes) != 0 && len(hashes) != len(urls) {
		return nil, fmt.Errorf("manifest has %d URLs but %d hashes", len(urls), len(hashes))
	}
	if len(subDirs) != 0 && len(subDirs) != len(urls) {
		return nil, fmt.Errorf("manifest has %d URLs but %d extract_dirs", len(urls), len(subDirs))
	}

	entries := make([]DownloadEntry, len(urls))
	for i, u := range urls {
		entry := DownloadEntry{URL: u}
		if len(hashes) != 0 {
			entry.Hash = hashes[i]
		}
		if len(subDirs) != 0 {
			entry.SubDir = subDirs[i]
		}
		entries[i] = entry
	}

	return entries, nil
}

// TODO: handle more bin formats (e.g. array of arrays)

func GetManifestBin(manifestPath string) ([]string, error) {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest json: %w", err)
	}

	return extractStringOrArray(jsonData["bin"]), nil
}

func GetManifestVersion(manifestPath string) (string, error) {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return "", fmt.Errorf("failed to get manifest json: %w", err)
	}

	version, ok := jsonData["version"].(string)
	if !ok {
		return "", fmt.Errorf("unable to find manifest version for %s", manifestPath)
	}

	return version, nil
}

func getManifestJson(manifestPath string) (map[string]any, error) {
	rawData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %q: %w", manifestPath, err)
	}

	var jsonData map[string]any
	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %q: %w", manifestPath, err)
	}

	return jsonData, nil
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
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}
