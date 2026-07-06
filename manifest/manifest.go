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

	var urls []string
	if urlVal, err := getArchSpecificProperty("url", architecture, jsonData); err == nil {
		urls = extractStringOrArray(urlVal)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("unable to find download URL for %s (arch %q)", manifestPath, architecture)
	}

	var hashes []string
	if hashVal, err := getArchSpecificProperty("hash", architecture, jsonData); err == nil {
		hashes = extractStringOrArray(hashVal)
	}

	var subDirs []string
	if subDirVal, err := getArchSpecificProperty("extract_dir", architecture, jsonData); err == nil {
		subDirs = extractStringOrArray(subDirVal)
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

type BinEntry struct {
	Exe   string
	Alias string
	Args  string
}

func GetManifestBin(manifestPath string, architecture string) ([]BinEntry, error) {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest json: %w", err)
	}

	var rawBinEntries [][]string
	if bins, err := getArchSpecificProperty("bin", architecture, jsonData); err == nil {
		rawBinEntries = extractStringOrArrayOrArrayOfArray(bins)
	}

	binEntries := make([]BinEntry, len(rawBinEntries))
	for i, rawBinEntry := range rawBinEntries {
		binEntry := BinEntry{}

		if len(rawBinEntry) > 0 {
			binEntry.Exe = rawBinEntry[0]
			binEntry.Alias = rawBinEntry[0]
		}
		if len(rawBinEntry) > 1 {
			binEntry.Alias = rawBinEntry[1]
		}
		if len(rawBinEntry) > 2 {
			binEntry.Args = rawBinEntry[2]
		}

		binEntries[i] = binEntry
	}

	return binEntries, nil
}

type ShortcutEntry struct {
	Exe  string
	Name string
	Args string
	Icon string
}

func GetManifestShortcut(manifestPath string, architecture string) ([]ShortcutEntry, error) {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest json: %w", err)
	}

	var rawShortcutEntries [][]string
	if shortcuts, err := getArchSpecificProperty("shortcuts", architecture, jsonData); err == nil {
		rawShortcutEntries = extractStringOrArrayOrArrayOfArray(shortcuts)
	}

	shortcutEntries := make([]ShortcutEntry, len(rawShortcutEntries))
	for i, rawShortcutEntry := range rawShortcutEntries {
		shortcutEntry := ShortcutEntry{}

		if len(rawShortcutEntry) < 2 {
			continue
		}

		shortcutEntry.Exe = rawShortcutEntry[0]
		shortcutEntry.Name = rawShortcutEntry[1]

		if len(rawShortcutEntry) > 2 {
			shortcutEntry.Args = rawShortcutEntry[2]
		}
		if len(rawShortcutEntry) > 3 {
			shortcutEntry.Icon = rawShortcutEntry[3]
		}
		shortcutEntries[i] = shortcutEntry
	}

	return shortcutEntries, nil
}

func GetManifestInnoSetup(manifestPath string) bool {
	jsonData, err := getManifestJson(manifestPath)
	if err != nil {
		return false
	}

	if val, ok := jsonData["innosetup"].(bool); ok {
		return val
	}

	return false
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

func getArchSpecificProperty(property string, architecture string, jsonData map[string]any) (any, error) {
	if archMap, ok := jsonData["architecture"].(map[string]any); ok {
		if archBlock, ok := archMap[architecture].(map[string]any); ok {
			if val, ok := archBlock[property]; ok {
				return val, nil
			}
		}
	}

	if val, ok := jsonData[property]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("failed to get %q from manifest json", property)
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

func extractStringOrArrayOrArrayOfArray(v any) [][]string {
	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		return [][]string{{val}}
	case []any:
		out := make([][]string, 0, len(val))
		for _, item := range val {
			if s := extractStringOrArray(item); len(s) > 0 {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}
