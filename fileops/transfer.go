package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func MergeDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		dstPath := filepath.Join(dst, relativePath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		if err := os.Rename(path, dstPath); err != nil {
			return fmt.Errorf("failed to move %s to %s: %w", path, dstPath, err)
		}
		return nil
	})
}

func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", src, err)
	}
	defer srcFile.Close()

	targetFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", dst, err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file %s: %w", dst, err)
	}

	if err := targetFile.Sync(); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}

	return nil
}
