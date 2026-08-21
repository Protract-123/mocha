package shim

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/Protract-123/mocha/manifest"
)

type Release struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	Version     string `json:"version"`
	PublishedAt string `json:"published_at"`
}

func GetLatestShimReleases() ([]Release, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/ScoopInstaller/Shim/releases", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("User-Agent", "Mocha-CLI")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed: %s", resp.Status)
	}

	var releases []struct {
		Name        string `json:"name"`
		TagName     string `json:"tag_name"`
		PublishedAt string `json:"published_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	latestReleases := map[string]Release{}

	for _, release := range releases {
		parts := strings.Split(release.TagName, "/")
		if len(parts) != 2 {
			continue
		}

		lang := parts[0]
		version := strings.TrimPrefix(parts[1], "v")

		current, exists := latestReleases[lang]
		if !exists || manifest.CompareVersions(current.Version, version) > 0 {
			latestReleases[lang] = Release{
				Name:        release.Name,
				Language:    lang,
				Version:     version,
				PublishedAt: release.PublishedAt,
			}
		}
	}

	result := slices.SortedFunc(maps.Values(latestReleases), func(a, b Release) int {
		return strings.Compare(a.Language, b.Language)
	})

	return result, nil
}
