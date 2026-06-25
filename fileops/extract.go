package fileops

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TODO: consider more secure extractors: e.g. https://github.com/hashicorp/go-extract or force 7z to be installed

func ExtractFile(filePath string, outputDir string) error {
	extension := filepath.Ext(filePath)
	switch extension {
	case ".zip":
		err := extractZip(filePath, outputDir)
		if err != nil {
			return fmt.Errorf("failed to extract %s to %s: %w", filePath, outputDir, err)
		}
	}

	return nil
}

func extractZip(zipPath string, outputDir string) error {
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
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}

	return nil
}

// TODO: handle zip bombs

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

	_, err = io.Copy(outFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy data to file %s: %w", file.Name, err)
	}

	err = outFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", file.Name, err)
	}

	return nil
}
