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

func ExtractInnoSetup(exePath string, outputDir string, subDir string) error {
	if _, err := exec.LookPath("innounp.exe"); err != nil {
		return fmt.Errorf("innounp is required to extract an InnoSetup")
	}

	var extractDir string

	if subDir == "" {
		extractDir = "{app}"
	} else if strings.HasPrefix(subDir, "{") {
		extractDir = subDir
	} else {
		extractDir = `{app}\` + subDir
	}

	cmd := exec.Command("innounp.exe", "-x", "-d"+outputDir, exePath, "-y", "-c"+extractDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract .exe: %w", err)
	}

	return nil
}

func Extract7z(archivePath string, outputDir string) error {
	if _, err := exec.LookPath("7z.exe"); err != nil {
		return fmt.Errorf("7z is required to extract with 7zip")
	}

	hasTar := hasTarArchive(archivePath)

	cmd := exec.Command("7z.exe", "x", archivePath, "-o"+outputDir, "-xr!*.nsis", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract %s: %w", filepath.Base(archivePath), err)
	}

	if hasTar {
		tarFile, err := findTarFile(outputDir)
		if err != nil {
			return fmt.Errorf("failed to find inner tar in %s: %w", filepath.Base(archivePath), err)
		}
		if err := Extract7z(tarFile, outputDir); err != nil {
			return fmt.Errorf("failed to extract inner tar: %w", err)
		}
		if err := os.Remove(tarFile); err != nil {
			return fmt.Errorf("failed to remove inner tar %s: %w", filepath.Base(tarFile), err)
		}
	}

	return nil
}

func hasTarArchive(path string) bool {
	strippedPath := strings.TrimSuffix(path, filepath.Ext(path))
	if strings.HasSuffix(strings.ToLower(strippedPath), ".tar") {
		return true
	}

	lowerPath := strings.ToLower(path)

	for _, suffix := range []string{".taz", ".tbz", ".tbz2", ".tgz", ".tpz", ".txz"} {
		if strings.HasSuffix(lowerPath, suffix) {
			return true
		}
	}

	return false
}

func findTarFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".tar") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no .tar file found in %s", dir)
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
