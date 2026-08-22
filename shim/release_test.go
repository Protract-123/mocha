package shim

import (
	"strings"
	"testing"

	"github.com/Protract-123/mocha/manifest"
)

func TestGetLatestShimReleases(tester *testing.T) {
	results, err := GetLatestShimReleases()
	if err != nil {
		tester.Fatalf("GetLatestShimReleases failed: %v", err)
	}

	prevLang := ""
	for _, release := range results {
		if prevLang != "" && strings.Compare(prevLang, release.Language) == 1 {
			tester.Errorf("release list is not sorted by language")
		}

		prevLang = release.Language
	}
}

func TestGetShimReleases(tester *testing.T) {
	results, err := GetShimReleases()
	if err != nil {
		tester.Fatalf("GetShimReleases failed: %v", err)
	}

	prevLang := ""
	prevVersion := ""

	for _, release := range results {
		if prevLang == "" {
			prevLang = release.Language
			continue
		}

		langComparison := strings.Compare(prevLang, release.Language)
		if langComparison == 1 {
			tester.Errorf("release list is not sorted by language")
		} else if langComparison == -1 {
			prevVersion = ""
		}

		if prevVersion == "" {
			prevVersion = release.Version
			prevLang = release.Language
			continue
		}

		versionComparison := manifest.CompareVersions(prevVersion, release.Version)
		if versionComparison == -1 {
			tester.Errorf("release list is not sorted by version")
		}

		prevLang = release.Language
		prevVersion = release.Version
	}
}

func TestFetchShimReleases(tester *testing.T) {
	results, err := fetchShimReleases()
	if err != nil {
		tester.Fatalf("fetchShimReleases failed: %v", err)
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
	}

	for _, lang := range requiredLanguages {
		if !seenLanguages[lang] {
			tester.Errorf("language %s not found in results", lang)
		}
	}
}
