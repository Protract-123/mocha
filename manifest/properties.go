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

func GetDownloadEntries(jsonData map[string]any, architecture string) ([]DownloadEntry, error) {
	var urls []string
	if val, err := getArchSpecificProperty("url", architecture, jsonData); err == nil {
		urls = extractAsArray(val)
	}

	var hashes []string
	if val, err := getArchSpecificProperty("hash", architecture, jsonData); err == nil {
		hashes = extractAsArray(val)
	}

	var subDirs []string
	if val, err := getArchSpecificProperty("extract_dir", architecture, jsonData); err == nil {
		subDirs = extractAsArray(val)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("unable to find download URL (arch %q)", architecture)
	}
	if len(hashes) != 0 && len(hashes) != len(urls) {
		return nil, fmt.Errorf("manifest has %d URLs but %d hashes", len(urls), len(hashes))
	}
	if len(subDirs) != 0 && len(subDirs) != len(urls) {
		return nil, fmt.Errorf("manifest has %d URLs but %d extract_dirs", len(urls), len(subDirs))
	}

	entries := make([]DownloadEntry, len(urls))
	for i, url := range urls {
		entry := DownloadEntry{URL: url}
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

type ExecutableEntry struct {
	Exe   string
	Alias string
	Args  string
}

func GetExecutableEntries(jsonData map[string]any, architecture string) ([]ExecutableEntry, error) {
	var rawEntries [][]string
	if val, err := getArchSpecificProperty("bin", architecture, jsonData); err == nil {
		rawEntries = extractAsArrayOfArray(val)
	}

	entries := make([]ExecutableEntry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := ExecutableEntry{}

		if len(rawEntry) > 0 {
			entry.Exe = rawEntry[0]
			entry.Alias = rawEntry[0]
		}
		if len(rawEntry) > 1 {
			entry.Alias = rawEntry[1]
		}
		if len(rawEntry) > 2 {
			entry.Args = rawEntry[2]
		}

		entries[i] = entry
	}

	return entries, nil
}

type ShortcutEntry struct {
	Exe  string
	Name string
	Args string
	Icon string
}

func GetShortcutEntries(jsonData map[string]any, architecture string) []ShortcutEntry {
	var rawEntries [][]string
	if val, err := getArchSpecificProperty("shortcuts", architecture, jsonData); err == nil {
		rawEntries = extractAsArrayOfArray(val)
	}

	entries := make([]ShortcutEntry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := ShortcutEntry{}

		if len(rawEntry) < 2 {
			continue
		}

		entry.Exe = rawEntry[0]
		entry.Name = rawEntry[1]

		if len(rawEntry) > 2 {
			entry.Args = rawEntry[2]
		}
		if len(rawEntry) > 3 {
			entry.Icon = rawEntry[3]
		}
		entries[i] = entry
	}

	return entries
}

type PersistEntry struct {
	Target string
	Source string
}

func GetPersistEntries(jsonData map[string]any) []PersistEntry {
	var rawEntries [][]string
	if val, ok := jsonData["persist"]; ok {
		rawEntries = extractAsArrayOfArray(val)
	}

	entries := make([]PersistEntry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry := PersistEntry{}

		if len(rawEntry) > 0 {
			entry.Target = rawEntry[0]
			entry.Source = rawEntry[0]
		}
		if len(rawEntry) > 1 {
			entry.Source = rawEntry[1]
		}

		entries[i] = entry
	}

	return entries
}

func GetInnoSetup(jsonData map[string]any) bool {
	if val, ok := jsonData["innosetup"].(bool); ok {
		return val
	}
	return false
}

func GetVersion(jsonData map[string]any) (string, error) {
	if val, ok := jsonData["version"].(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("failed to get version from json")
}

func GetJson(manifestPath string) (map[string]any, error) {
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

func extractAsArray(v any) []string {
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

func extractAsArrayOfArray(v any) [][]string {
	switch val := v.(type) {
	case string:
		if val == "" {
			return nil
		}
		return [][]string{{val}}
	case []any:
		out := make([][]string, 0, len(val))
		for _, item := range val {
			if s := extractAsArray(item); len(s) > 0 {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}
