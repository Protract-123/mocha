package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Info struct {
	Name   string
	Target string
}

func CreateShim(name string, path string, mochaDir string) error {
	shimDirectory := filepath.Join(mochaDir, "shims")
	if err := os.MkdirAll(shimDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create shim directory: %w", err)
	}

	fileExtension := filepath.Ext(path)

	if fileExtension == ".exe" || fileExtension == ".com" {
		if err := CreateExeShim(name, path, mochaDir); err != nil {
			return fmt.Errorf("failed to create exe shim: %w", err)
		}
	}

	return nil
}

func DeleteShim(name string, mochaDir string) error {
	shimsDir := filepath.Join(mochaDir, "shims")

	shims, err := os.ReadDir(shimsDir)
	if err != nil {
		return fmt.Errorf("failed to read shims directory: %w", err)
	}

	deletion := false

	for _, shim := range shims {
		if shim.IsDir() {
			continue
		}

		shimName := strings.TrimSuffix(shim.Name(), filepath.Ext(shim.Name()))

		if shimName == name {
			if err := os.Remove(filepath.Join(shimsDir, shim.Name())); err != nil {
				return fmt.Errorf("failed to remove shim file %s: %w", shim.Name(), err)
			}
			deletion = true
		}
	}

	if !deletion {
		return fmt.Errorf("no shim found for %q", name)
	}

	return nil
}

func GetAllShims(mochaDir string) ([]Info, error) {
	shimsDir := filepath.Join(mochaDir, "shims")

	shims, err := os.ReadDir(shimsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read shims directory: %w", err)
	}

	var shimInfo []Info

	for _, shim := range shims {
		if shim.IsDir() {
			continue
		}

		if filepath.Ext(shim.Name()) == ".shim" {
			path := filepath.Join(shimsDir, shim.Name())

			shimBytes, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read shim file %s: %w", shim.Name(), err)
			}

			name := strings.TrimSuffix(shim.Name(), ".shim")
			target := strings.Split(string(shimBytes), "=")[1]

			shimInfo = append(shimInfo, Info{name, target})
		}
	}

	return shimInfo, nil
}
