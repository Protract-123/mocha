package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Ref struct {
	Bucket  string
	Name    string
	Version string
}

func ParseAppString(appString string) (Ref, error) {
	if appString == "" {
		return Ref{}, fmt.Errorf("app string is empty")
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
		return Ref{}, fmt.Errorf("app name not parsed correctly")
	} else if appRef.Bucket == "" && strings.Contains(appString, "/") {
		return Ref{}, fmt.Errorf("app bucket name not parsed correctly")
	} else if appRef.Version == "" && strings.Contains(appString, "@") {
		return Ref{}, fmt.Errorf("app version not parsed correctly")
	}

	return appRef, nil
}

func PopulateAppRef(appRef Ref, mochaDir string) (Ref, error) {
	if appRef.Name == "" {
		return appRef, fmt.Errorf("app name is required to populate the app ref")
	}

	if appRef.Bucket == "" {
		bucketsDir := filepath.Join(mochaDir, "buckets")
		buckets, err := os.ReadDir(bucketsDir)
		if err != nil {
			return appRef, err
		}

		for _, dirEntry := range buckets {
			if !dirEntry.IsDir() {
				continue
			}

			path, err := GetManifestPath(dirEntry.Name(), appRef.Name, mochaDir)
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return appRef, err
			}

			if path != "" {
				appRef.Bucket = dirEntry.Name()
				break
			}
		}

		if appRef.Bucket == "" {
			return appRef, fmt.Errorf("could not find app %q in any bucket", appRef.Name)
		}
	}

	if appRef.Version == "" {
		manifestPath, err := GetManifestPath(appRef.Bucket, appRef.Name, mochaDir)
		if err != nil {
			return appRef, err
		}

		version, err := GetManifestVersion(manifestPath)
		if err != nil {
			return appRef, err
		}

		appRef.Version = version
	}

	return appRef, nil
}
