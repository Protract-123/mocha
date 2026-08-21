package shim

import (
	"testing"
)

func TestGetLatestShimReleases(tester *testing.T) {
	results, err := GetLatestShimReleases()
	if err != nil {
		tester.Errorf("getLatestShimReleases failed: %v", err)
	}

	requiredLanguages := []string{
		"zig",
		"cs",
		"rust",
		"cpp",
	}

	seenLanguages := map[string]bool{}

	for _, release := range results {
		seenLanguages[release.Language] = true
		tester.Logf("name: %s, language: %s, version: %s", release.Name, release.Language, release.Version)
	}

	for _, lang := range requiredLanguages {
		if !seenLanguages[lang] {
			tester.Errorf("language %s not found in results", lang)
		}
	}
}
