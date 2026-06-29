package shim

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Protract-123/mocha/fileops"
	"github.com/Protract-123/mocha/output"
)

var shimSources = []string{
	"https://github.com/ScoopInstaller/Shim/releases/download/rust%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/zig%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/cs%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/cpp%2Fv0.1.0/",
}

func InitShimBinary(mochaDir string) error {
	shimExe := filepath.Join(mochaDir, "shim.exe")

	if _, err := os.Stat(shimExe); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check if shim.exe exists: %w", err)
	}

	arch, err := getCPUArch()
	if err != nil {
		return fmt.Errorf("failed to get cpu architecture: %w", err)
	}

	tempDirectory := filepath.Join(mochaDir, "temp")
	if err := os.MkdirAll(tempDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDirectory)

	zipName := fmt.Sprintf("shim-%s.zip", arch)
	zipPath := filepath.Join(tempDirectory, zipName)

	downloadURL, err := url.JoinPath(shimSources[3], zipName)
	if err != nil {
		return fmt.Errorf("failed to create download url: %w", err)
	}

	output.LogOutput(fmt.Sprintf("Downloading %s from %s\n", zipName, downloadURL))

	if err := fileops.DownloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("failed to download %s: %w", zipName, err)
	}

	if err := fileops.ExtractZip(zipPath, tempDirectory); err != nil {
		return fmt.Errorf("failed to extract %s: %w", zipName, err)
	}

	shimExeBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe: %w", err)
	}

	checksumBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe.sha256"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe.sha256: %w", err)
	}

	checksum := strings.Split(strings.TrimSpace(string(checksumBytes)), " ")[0]
	sumBytes := sha256.Sum256(shimExeBytes)
	sum := hex.EncodeToString(sumBytes[:])

	if sum != checksum {
		return fmt.Errorf("shim.exe hash does not match shim.exe.sha256")
	}

	if err := os.WriteFile(shimExe, shimExeBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write shim.exe to %s: %w", shimExe, err)
	}

	return nil
}

func getCPUArch() (string, error) {
	switch runtime.GOARCH {
	case "386":
		return "x86", nil
	case "amd64":
		return "x64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("cpu architecture %q is unsupported", runtime.GOARCH)
	}
}
