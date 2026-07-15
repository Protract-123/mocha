package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/pkg"
)

type UninstallCommand struct {
	Apps []string `arg:"positional,required" help:"apps to uninstall; omit @version to remove all installed versions (e.g. git, bat@0.26.1)"`
}

func (cmd *UninstallCommand) Run(mochaDir string) error {
	for _, refString := range cmd.Apps {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		appDir := filepath.Join(mochaDir, "apps", manifestRef.Name)
		var deletionDir string

		if manifestRef.Version != "" {
			deletionDir = filepath.Join(appDir, manifestRef.Version)
		} else {
			deletionDir = filepath.Join(appDir)
		}

		if _, err := os.Stat(deletionDir); os.IsNotExist(err) {
			return fmt.Errorf("%q is not installed", refString)
		} else if err != nil {
			return fmt.Errorf("failed to check if %q exists: %w", refString, err)
		}

		activeVersion, err := pkg.GetActiveVersion(manifestRef.Name, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get active version of %q: %w", manifestRef.Name, err)
		}

		if activeVersion != "" && (activeVersion == manifestRef.Version || manifestRef.Version == "") {
			if err := pkg.UnlinkApp(manifestRef.Name, mochaDir); err != nil {
				return fmt.Errorf("failed to unlink app %q: %w", manifestRef.Name, err)
			}
		}

		if err := os.RemoveAll(deletionDir); err != nil {
			return fmt.Errorf("failed to uninstall %q: %w", refString, err)
		}

		if deletionDir != appDir {
			files, err := os.ReadDir(appDir)
			if err != nil {
				return fmt.Errorf("failed to get files in %q: %w", deletionDir, err)
			}

			if len(files) == 0 {
				if err := os.RemoveAll(appDir); err != nil {
					return fmt.Errorf("failed to remove %q: %w", appDir, err)
				}
			}
		}
	}

	return nil
}
