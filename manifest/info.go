package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Info struct {
	Name         string
	Bucket       string
	Version      string
	ManifestPath string
}

func ParseSpec(spec string) (Info, error) {
	malformedSpecError := fmt.Errorf("invalid manifest spec %q, expected format %q", spec, "[bucket/]manifest[@version]")

	if spec == "" || strings.Count(spec, "@") > 1 || strings.Count(spec, "/") > 1 {
		return Info{}, malformedSpecError
	}

	info := Info{}
	unparsedPortion := spec

	if bucket, rest, found := strings.Cut(unparsedPortion, "/"); found {
		info.Bucket = bucket
		unparsedPortion = rest
	}

	if name, version, found := strings.Cut(unparsedPortion, "@"); found {
		info.Name = name
		info.Version = version
	} else {
		info.Name = unparsedPortion
	}

	if info.Name == "" {
		return Info{}, malformedSpecError
	} else if info.Bucket == "" && strings.Contains(spec, "/") {
		return Info{}, malformedSpecError
	} else if info.Version == "" && strings.Contains(spec, "@") {
		return Info{}, malformedSpecError
	}

	return info, nil
}

func PopulateInfo(info Info, mochaDir string) (Info, error) {
	if info.Name == "" {
		return Info{}, fmt.Errorf("manifest name is empty")
	}

	if info.Bucket == "" {
		buckets, err := os.ReadDir(filepath.Join(mochaDir, "buckets"))
		if err != nil {
			return Info{}, fmt.Errorf("failed to get all buckets: %w", err)
		}

		for _, bucket := range buckets {
			if !bucket.IsDir() {
				continue
			}

			manifestPath := filepath.Join(mochaDir, "buckets", bucket.Name(), "bucket", fmt.Sprintf("%s.json", info.Name))
			if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
				continue
			} else if err != nil {
				return Info{}, fmt.Errorf("failed to confirm if manifest exists at %q: %w", manifestPath, err)
			}

			info.Bucket = bucket.Name()
			break
		}

		if info.Bucket == "" {
			return Info{}, fmt.Errorf("failed to find manifest %q in buckets", info.Name)
		}
	}

	manifestPath := filepath.Join(mochaDir, "buckets", info.Bucket, "bucket", fmt.Sprintf("%s.json", info.Name))
	if _, err := os.Stat(manifestPath); err != nil {
		return Info{}, fmt.Errorf("failed to find manifest %q in bucket %q: %w", info.Name, info.Bucket, err)
	}

	info.ManifestPath = manifestPath

	if info.Version == "" {
		manifestJson, err := GetJson(manifestPath)
		if err != nil {
			return Info{}, fmt.Errorf("failed to get manifest JSON: %w", err)
		}

		version, err := GetVersion(manifestJson)
		if err != nil {
			return Info{}, fmt.Errorf("failed to get manifest version for %q in bucket %q: %w", info.Name, info.Bucket, err)
		}

		info.Version = version
	}

	return info, nil
}
