package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Info struct {
	Name        string
	Destination string
}

func CreateShim(name string, path string, mochaDir string) error {
	fileExtension := filepath.Ext(path)

	if fileExtension == ".exe" || fileExtension == ".com" {
		err := CreateExeShim(name, path, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to create exe shim: %w", err)
		}
	}

	return nil
}

func DeleteShim(name string, mochaDir string) error {
	shimsDir := filepath.Join(mochaDir, "shims")

	files, err := os.ReadDir(shimsDir)
	if err != nil {
		return fmt.Errorf("failed to read shims directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileExtension := filepath.Ext(file.Name())
		fileName := strings.TrimSuffix(file.Name(), fileExtension)

		if fileName == name {
			err := os.Remove(filepath.Join(shimsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to remove shim file %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}

func GetAllShims(mochaDir string) ([]Info, error) {
	shimsDir := filepath.Join(mochaDir, "shims")

	files, err := os.ReadDir(shimsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read shims directory: %w", err)
	}

	var shims []Info

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".shim" {
			path := filepath.Join(shimsDir, file.Name())

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read shim file %s: %w", file.Name(), err)
			}

			name := strings.TrimSuffix(file.Name(), ".shim")
			destination := strings.Split(string(data), "=")[1]

			shims = append(shims, Info{name, destination})
		}
	}

	return shims, nil
}
