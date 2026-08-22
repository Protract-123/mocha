package shim

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Protract-123/mocha/manifest"
)

type Release struct {
	Name        string    `json:"name"`
	Language    string    `json:"language"`
	Version     string    `json:"version"`
	PublishedAt time.Time `json:"published_at"`
}

func GetLatestShimReleases() ([]Release, error) {
	releases, err := fetchShimReleases()
	if err != nil {
		return nil, err
	}

	latestReleases := map[string]Release{}

	for _, release := range releases {
		current, exists := latestReleases[release.Language]
		if !exists || manifest.CompareVersions(current.Version, release.Version) > 0 {
			latestReleases[release.Language] = release
		}
	}

	result := slices.SortedFunc(maps.Values(latestReleases), func(a, b Release) int {
		return strings.Compare(a.Language, b.Language)
	})

	return result, nil
}

func GetShimReleases() ([]Release, error) {
	releases, err := fetchShimReleases()
	if err != nil {
		return nil, err
	}

	slices.SortFunc(releases, func(a, b Release) int {
		if dateCmp := b.PublishedAt.Compare(a.PublishedAt); dateCmp != 0 {
			return dateCmp
		}
		if langCmp := strings.Compare(a.Language, b.Language); langCmp != 0 {
			return langCmp
		}
		return manifest.CompareVersions(b.Version, a.Version)
	})

	return releases, nil
}

func fetchShimReleases() ([]Release, error) {
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

	shimReleases := make([]Release, 0, len(releases))

	for _, release := range releases {
		parts := strings.Split(release.TagName, "/")
		if len(parts) != 2 {
			continue
		}

		lang := parts[0]
		version := strings.TrimPrefix(parts[1], "v")

		publishAt, err := time.Parse(time.RFC3339, release.PublishedAt)
		if err != nil {
			continue
		}

		shimReleases = append(shimReleases, Release{
			Name:        release.Name,
			Language:    lang,
			Version:     version,
			PublishedAt: publishAt,
		})
	}

	return shimReleases, nil
}
