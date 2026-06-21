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
	shimDir := filepath.Join(mochaDir, "shims")
	shimPath := filepath.Join(mochaDir, "shim.exe")

	err := os.MkdirAll(shimDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Stat(shimPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		return nil
	}

	arch, err := app.GetSystemArch()
	if err != nil {
		return err
	}

	if arch == "64bit" {
		arch = "x64"
	} else if arch == "32bit" {
		arch = "x86"
	}

	fileName := fmt.Sprintf("shim-%s.zip", arch)
	outputPath := filepath.Join(mochaDir, fileName)

	downloadUrl, err := url.JoinPath(shimSources[3], fileName)
	if err != nil {
		return err
	}
	fmt.Printf("Downloading shim from %s\n", downloadUrl)

	err = fileops.DownloadFile(downloadUrl, outputPath)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			println(err.Error())
		}
	}(outputPath)

	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			println(err.Error())
		}
	}(reader)

	var shimData []byte
	var hashString string

	for _, file := range reader.File {
		if file.Name == "shim.exe" {
			fileHandle, err := file.Open()
			if err != nil {
				return err
			}

			shimData, err = io.ReadAll(fileHandle)
			if err != nil {
				return err
			}

			err = fileHandle.Close()
			if err != nil {
				return err
			}
		}

		if file.Name == "shim.exe.sha256" {
			fileHandle, err := file.Open()
			if err != nil {
				return err
			}

			content, err := io.ReadAll(fileHandle)
			if err != nil {
				return err
			}

			fileContent := strings.TrimSpace(string(content))
			hashString = strings.Split(fileContent, " ")[0]

			err = fileHandle.Close()
			if err != nil {
				return err
			}
		}
	}

	if shimData == nil {
		return fmt.Errorf("shim.exe not found in archive")
	}

	if hashString == "" {
		return fmt.Errorf("shim.exe.sha256 not found in archive")
	}

	sum := sha256.Sum256(shimData)
	fileHash := hex.EncodeToString(sum[:])
	if fileHash != hashString {
		return fmt.Errorf("file hash mismatch")
	}

	return os.WriteFile(shimPath, shimData, 0755)
}
