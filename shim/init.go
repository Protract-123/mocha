package shim

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/app"
	"github.com/Protract-123/mocha/fileops"
)

var shimSources = []string{
	"https://github.com/ScoopInstaller/Shim/releases/download/rust%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/zig%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/cs%2Fv0.1.0/",
	"https://github.com/ScoopInstaller/Shim/releases/download/cpp%2Fv0.1.0/",
}

func InitShims(mochaDir string) error {
	tmpDir := filepath.Join(mochaDir, "tmp")

	shimDir := filepath.Join(mochaDir, "shims")
	shimPath := filepath.Join(mochaDir, "shim.exe")

	err := os.MkdirAll(shimDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create shim directory: %w", err)
	}

	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create tmp directory: %w", err)
	}

	_, err = os.Stat(shimPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if shim.exe exists: %w", err)
	} else if err == nil {
		return nil
	}

	arch, err := app.GetSystemArch()
	if err != nil {
		return fmt.Errorf("failed to get cpu architecture: %w", err)
	}

	if arch == "64bit" {
		arch = "x64"
	} else if arch == "32bit" {
		arch = "x86"
	}

	fileName := fmt.Sprintf("shim-%s.zip", arch)
	outputPath := filepath.Join(tmpDir, fileName)

	downloadUrl, err := url.JoinPath(shimSources[3], fileName)
	if err != nil {
		return fmt.Errorf("failed to create download url: %w", err)
	}

	fmt.Printf("Downloading %s from %s\n", fileName, downloadUrl)

	err = fileops.DownloadFile(downloadUrl, outputPath)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", fileName, err)
	}
	defer os.Remove(outputPath)

	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	var shimData []byte
	var hashString string

	for _, file := range reader.File {
		if file.Name == "shim.exe" {
			shimData, err = readZipFile(file)
			if err != nil {
				return err
			}
		}

		if file.Name == "shim.exe.sha256" {
			content, err := readZipFile(file)
			if err != nil {
				return err
			}
			hashString = strings.Split(strings.TrimSpace(string(content)), " ")[0]
		}
	}

	if shimData == nil {
		return fmt.Errorf("shim.exe not found in %s", fileName)
	}

	if hashString == "" {
		return fmt.Errorf("shim.exe.sha256 not found in %s", fileName)
	}

	sum := sha256.Sum256(shimData)
	fileHash := hex.EncodeToString(sum[:])
	if fileHash != hashString {
		return fmt.Errorf("shim.exe hash does not match shim.exe.sha256")
	}

	err = os.WriteFile(shimPath, shimData, 0755)
	if err != nil {
		return fmt.Errorf("failed to write shim.exe to %s: %w", shimPath, err)
	}

	return nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", file.Name, err)
	}
	defer fileHandle.Close()

	data, err := io.ReadAll(fileHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", file.Name, err)
	}

	return data, nil
}
