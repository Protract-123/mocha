package shim

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/fileops"
)

func InstallBinary(release Release, arch string, mochaDir string) error {
	tempDirectory := filepath.Join(mochaDir, "temp")
	if err := os.MkdirAll(tempDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDirectory)

	zipName := fmt.Sprintf("shim-%s.zip", arch)
	zipPath := filepath.Join(tempDirectory, zipName)

	downloadURL, err := url.JoinPath("https://github.com/ScoopInstaller/Shim/releases/download", release.Language, "v"+release.Version, zipName)
	if err != nil {
		return fmt.Errorf("failed to create download url: %w", err)
	}

	if err := fileops.DownloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("failed to download %s: %w", zipName, err)
	}

	if err := fileops.ExtractZip(zipPath, tempDirectory); err != nil {
		return fmt.Errorf("failed to extract %s: %w", zipName, err)
	}

	binaryBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe: %w", err)
	}

	checksumBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe.sha256"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe.sha256: %w", err)
	}

	sumBytes := sha256.Sum256(binaryBytes)
	sum := hex.EncodeToString(sumBytes[:])

	checksum := strings.Split(strings.TrimSpace(string(checksumBytes)), " ")[0]

	if sum != checksum {
		return fmt.Errorf("shim.exe hash does not match shim.exe.sha256")
	}

	binaryPath := filepath.Join(mochaDir, "shim.exe")
	if err := os.WriteFile(binaryPath, binaryBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write shim.exe to %s: %w", binaryPath, err)
	}

	return nil
}
