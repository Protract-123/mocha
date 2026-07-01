package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/shim"
)

type UninstallCommand struct {
	ManifestReferences []string `arg:"positional,required"`
}

func (cmd *UninstallCommand) Run(mochaDir string) error {
	for _, refString := range cmd.ManifestReferences {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		appDir := filepath.Join(mochaDir, "apps")
		var deletionDir string

		if manifestRef.Version != "" {
			deletionDir = filepath.Join(appDir, manifestRef.Name, manifestRef.Version)
		} else {
			deletionDir = filepath.Join(appDir, manifestRef.Name)
		}

		if _, err := os.Stat(deletionDir); os.IsNotExist(err) {
			return fmt.Errorf("%q is not installed", refString)
		} else if err != nil {
			return fmt.Errorf("failed to check if %q exists: %w", refString, err)
		}

		if err := os.RemoveAll(deletionDir); err != nil {
			return fmt.Errorf("failed to uninstall %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to populate manifest ref %q: %w", refString, err)
		}

		binaries, err := manifest.GetManifestBin(manifestRef.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get shims to remove: %w", err)
		}

		for _, binary := range binaries {
			shimName := strings.TrimSuffix(filepath.Base(binary), filepath.Ext(binary))
			if err := shim.DeleteShim(shimName, mochaDir); err != nil {
				return fmt.Errorf("failed to remove shim %q: %w", shimName, err)
			}
		}
	}

	return nil
}
