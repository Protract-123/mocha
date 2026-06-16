package manifest

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetManifestPath(bucketName string, manifestName string, mochaDir string) (string, error) {
	manifestPath := filepath.Join(mochaDir, "buckets", bucketName, "bucket", fmt.Sprintf("%s.json", manifestName))

	_, err := os.Stat(manifestPath)
	if err != nil {
		return "", err
	}

	return manifestPath, nil
}
