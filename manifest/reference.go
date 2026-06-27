package manifest

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

type BadRefError struct {
	providedRef    string
	expectedFormat string
}

func (error BadRefError) Error() string {
	return fmt.Sprintf("invalid manifest reference %q, expected %q", error.providedRef, error.expectedFormat)
}

func ParseRefString(refString string) (Ref, error) {
	if refString == "" {
		return Ref{}, BadRefError{refString, "[bucket/]manifest[@version]"}
	}

	manifestRef := Ref{}
	unparsedPortion := refString

	if bucket, rest, found := strings.Cut(unparsedPortion, "/"); found {
		manifestRef.Bucket = bucket
		unparsedPortion = rest
	}

	if name, version, found := strings.Cut(unparsedPortion, "@"); found {
		manifestRef.Name = name
		manifestRef.Version = version
	} else {
		manifestRef.Name = unparsedPortion
	}

	if manifestRef.Name == "" {
		return Ref{}, BadRefError{refString, "[bucket/]manifest[@version]"}
	} else if manifestRef.Bucket == "" && strings.Contains(refString, "/") {
		return Ref{}, BadRefError{refString, "bucket/manifest"}
	} else if manifestRef.Version == "" && strings.Contains(refString, "@") {
		return Ref{}, BadRefError{refString, "manifest@version"}
	}

	return manifestRef, nil
}

func PopulateRef(ref Ref, mochaDir string) (Ref, error) {
	if ref.Name == "" {
		return Ref{}, fmt.Errorf("manifest name is empty")
	}

	if ref.Bucket == "" {
		bucketsDir := filepath.Join(mochaDir, "buckets")
		buckets, err := os.ReadDir(bucketsDir)
		if err != nil {
			return Ref{}, fmt.Errorf("failed to get all buckets: %w", err)
		}

		for _, dirEntry := range buckets {
			if !dirEntry.IsDir() {
				continue
			}

			manifestPath := filepath.Join(mochaDir, "buckets", dirEntry.Name(), "bucket", fmt.Sprintf("%s.json", ref.Name))
			_, err := os.Stat(manifestPath)

			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return ref, fmt.Errorf("failed to confirm if manifest exists at %q: %w", manifestPath, err)
			}

			ref.Bucket = dirEntry.Name()
			break
		}

		if ref.Bucket == "" {
			return ref, fmt.Errorf("failed to find manifest %q in buckets", ref.Name)
		}
	}

	manifestPath := filepath.Join(mochaDir, "buckets", ref.Bucket, "bucket", fmt.Sprintf("%s.json", ref.Name))
	_, err := os.Stat(manifestPath)

	ref.ManifestPath = manifestPath

	if err != nil {
		return Ref{}, fmt.Errorf("failed to find manifest %q in bucket %q: %w", ref.Name, ref.Bucket, err)
	}

	if ref.Version == "" {
		version, err := GetManifestVersion(manifestPath)
		if err != nil {
			return ref, fmt.Errorf("failed to get manifest version for %q in bucket %q: %w", ref.Name, ref.Bucket, err)
		}

		ref.Version = version
	}

	return ref, nil
}
