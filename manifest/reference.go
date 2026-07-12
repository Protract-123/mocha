package manifest

import (
	"errors"
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

func ParseRefString(refString string) (Ref, error) {
	malformedRefError := fmt.Errorf("invalid manifest reference %q, expected format %q", refString, "[bucket/]manifest[@version]")

	if refString == "" || strings.Count(refString, "@") > 1 || strings.Count(refString, "/") > 1 {
		return Ref{}, malformedRefError
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
		return Ref{}, malformedRefError
	} else if manifestRef.Bucket == "" && strings.Contains(refString, "/") {
		return Ref{}, malformedRefError
	} else if manifestRef.Version == "" && strings.Contains(refString, "@") {
		return Ref{}, malformedRefError
	}

	return manifestRef, nil
}

func PopulateRef(ref Ref, mochaDir string) (Ref, error) {
	if ref.Name == "" {
		return Ref{}, fmt.Errorf("manifest name is empty")
	}

	if ref.Bucket == "" {
		buckets, err := os.ReadDir(filepath.Join(mochaDir, "buckets"))
		if err != nil {
			return Ref{}, fmt.Errorf("failed to get all buckets: %w", err)
		}

		for _, bucket := range buckets {
			if !bucket.IsDir() {
				continue
			}

			manifestPath := filepath.Join(mochaDir, "buckets", bucket.Name(), "bucket", fmt.Sprintf("%s.json", ref.Name))
			if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
				continue
			} else if err != nil {
				return Ref{}, fmt.Errorf("failed to confirm if manifest exists at %q: %w", manifestPath, err)
			}

			ref.Bucket = bucket.Name()
			break
		}

		if ref.Bucket == "" {
			return Ref{}, fmt.Errorf("failed to find manifest %q in buckets", ref.Name)
		}
	}

	manifestPath := filepath.Join(mochaDir, "buckets", ref.Bucket, "bucket", fmt.Sprintf("%s.json", ref.Name))
	if _, err := os.Stat(manifestPath); err != nil {
		return Ref{}, fmt.Errorf("failed to find manifest %q in bucket %q: %w", ref.Name, ref.Bucket, err)
	}

	ref.ManifestPath = manifestPath

	if ref.Version == "" {
		version, err := GetManifestVersion(manifestPath)
		if err != nil {
			return Ref{}, fmt.Errorf("failed to get manifest version for %q in bucket %q: %w", ref.Name, ref.Bucket, err)
		}

		ref.Version = version
	}

	return ref, nil
}
