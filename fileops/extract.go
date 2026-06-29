package fileops

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: consider more secure extractors: e.g. https://github.com/hashicorp/go-extract or force 7z to be installed

func ExtractMsi(msiPath string, outputDir string) error {
	cmd := exec.Command("msiexec.exe", "/a", msiPath, "/qn", fmt.Sprintf("TARGETDIR=%s", outputDir))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract msi: %w", err)
	}
	return nil
}

func ExtractZip(zipPath string, outputDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, file := range reader.File {
		err = extractZipFile(file, outputDir)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}

	return nil
}

// TODO: improve zip bomb protection

func extractZipFile(file *zip.File, outputDir string) error {
	outputPath := filepath.Join(outputDir, file.Name)

	relativePath, err := filepath.Rel(outputDir, outputPath)
	if err != nil || strings.HasPrefix(relativePath, "..") || strings.Contains(relativePath, ":") {
		return fmt.Errorf("illegal file path in zip: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputPath, err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputPath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	srcFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file.Name, err)
	}
	defer srcFile.Close()

	bytesWritten, err := io.Copy(outFile, io.LimitReader(srcFile, 512*1024*1024+1))
	if err != nil {
		return fmt.Errorf("failed to copy data to file %s: %w", file.Name, err)
	}
	if bytesWritten > 512*1024*1024 {
		return fmt.Errorf("file %s is too large", file.Name)
	}

	err = outFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", file.Name, err)
	}

	return nil
}
