package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Ref struct {
	Name         string
	Bucket       string
	Version      string
	ManifestPath string
}

type BadAppRefError struct {
	providedAppRef string
	expectedFormat string
}

func (error BadAppRefError) Error() string {
	return fmt.Sprintf("invalid app reference %q, expected %q", error.providedAppRef, error.expectedFormat)
}

func ParseAppString(appString string) (Ref, error) {
	if appString == "" {
		return Ref{}, BadAppRefError{appString, "[bucket/]app[@version]"}
	}

	appRef := Ref{}
	unparsedPortion := appString

	if bucket, rest, found := strings.Cut(unparsedPortion, "/"); found {
		appRef.Bucket = bucket
		unparsedPortion = rest
	}

	if name, version, found := strings.Cut(unparsedPortion, "@"); found {
		appRef.Name = name
		appRef.Version = version
	} else {
		appRef.Name = unparsedPortion
	}

	if appRef.Name == "" {
		return Ref{}, BadAppRefError{appString, "[bucket/]app[@version]"}
	} else if appRef.Bucket == "" && strings.Contains(appString, "/") {
		return Ref{}, BadAppRefError{appString, "bucket/app"}
	} else if appRef.Version == "" && strings.Contains(appString, "@") {
		return Ref{}, BadAppRefError{appString, "app@version"}
	}

	return appRef, nil
}

func PopulateAppRef(appRef Ref, mochaDir string) (Ref, error) {
	if appRef.Name == "" {
		return Ref{}, fmt.Errorf("app name is empty")
	}

	if appRef.Bucket == "" {
		bucketsDir := filepath.Join(mochaDir, "buckets")
		buckets, err := os.ReadDir(bucketsDir)
		if err != nil {
			return Ref{}, fmt.Errorf("failed to get all buckets: %w", err)
		}

		for _, dirEntry := range buckets {
			if !dirEntry.IsDir() {
				continue
			}

			manifestPath := filepath.Join(mochaDir, "buckets", dirEntry.Name(), "bucket", fmt.Sprintf("%s.json", appRef.Name))
			_, err := os.Stat(manifestPath)

			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return appRef, fmt.Errorf("failed to confirm if manifest exists at %q: %w", manifestPath, err)
			}

			appRef.Bucket = dirEntry.Name()
			break
		}

		if appRef.Bucket == "" {
			return appRef, fmt.Errorf("failed to find app %q in buckets", appRef.Name)
		}
	}

	manifestPath := filepath.Join(mochaDir, "buckets", appRef.Bucket, "bucket", fmt.Sprintf("%s.json", appRef.Name))
	_, err := os.Stat(manifestPath)

	appRef.ManifestPath = manifestPath

	if err != nil {
		return Ref{}, fmt.Errorf("failed to find app %q in bucket %q: %w", appRef.Name, appRef.Bucket, err)
	}

	if appRef.Version == "" {
		version, err := GetManifestVersion(manifestPath)
		if err != nil {
			return appRef, fmt.Errorf("failed to get app version for %q in bucket %q: %w", appRef.Name, appRef.Bucket, err)
		}

		appRef.Version = version
	}

	return appRef, nil
}
