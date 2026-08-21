package uninstall

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
)

type Command struct {
	Apps []string `arg:"positional,required" help:"apps to uninstall; omit @version to remove all installed versions (e.g. git, bat@0.26.1)"`
}

func (cmd *Command) Run() error {
	mochaDir := config.Current().MochaDirectory

	for _, spec := range cmd.Apps {
		info, err := manifest.ParseSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse manifest spec %q: %w", spec, err)
		}

		appDir := filepath.Join(mochaDir, "apps", info.Name)
		var deletionDir string

		if info.Version != "" {
			deletionDir = filepath.Join(appDir, info.Version)
		} else {
			deletionDir = filepath.Join(appDir)
		}

		if _, err := os.Stat(deletionDir); os.IsNotExist(err) {
			return fmt.Errorf("%q is not installed", spec)
		} else if err != nil {
			return fmt.Errorf("failed to check if %q exists: %w", spec, err)
		}

		activeVersion, err := pkg.GetActiveVersion(info.Name, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get active version of %q: %w", info.Name, err)
		}

		if activeVersion != "" && (activeVersion == info.Version || info.Version == "") {
			output.LogInfo("removing shims/shortcuts for %q", spec)
			if err := pkg.Unlink(info.Name, mochaDir); err != nil {
				return fmt.Errorf("failed to unlink app %q: %w", info.Name, err)
			}
		}

		output.LogInfo("uninstalling %q", spec)
		if err := os.RemoveAll(deletionDir); err != nil {
			return fmt.Errorf("failed to uninstall %q: %w", spec, err)
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

		output.LogSuccess("successfully uninstalled %q", spec)
	}

	return nil
}
