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
	downloadArch, err := pkg.GetDownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	for _, refString := range cmd.Apps {
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

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to populate manifest ref %q: %w", refString, err)
		}

		if err := pkg.UnlinkApp(manifestRef, downloadArch, mochaDir, deletionDir); err != nil {
			return fmt.Errorf("failed to unlink %q: %w", refString, err)
		}

		if err := os.RemoveAll(deletionDir); err != nil {
			return fmt.Errorf("failed to uninstall %q: %w", refString, err)
		}
	}

	return nil
}
