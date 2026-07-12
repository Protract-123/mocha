package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
	"github.com/alexflint/go-arg"
)

type UpgradeCommand struct {
	Apps []string `arg:"positional" help:"apps to upgrade (e.g. git, 7zip)"`
	All  bool     `arg:"-a,--all" help:"upgrade all apps"`
}

func (cmd *UpgradeCommand) Run(mochaDir string) error {
	if len(cmd.Apps) == 0 && !cmd.All {
		return arg.ErrHelp
	}

	downloadArch, err := pkg.GetDownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	var appList []string
	if cmd.All {
		apps, err := os.ReadDir(filepath.Join(mochaDir, "apps"))
		if err != nil {
			return fmt.Errorf("cannot read apps directory: %w", err)
		}

		for _, app := range apps {
			appList = append(appList, app.Name())
		}
	} else {
		appList = cmd.Apps
	}

	for _, app := range appList {
		infoPath := filepath.Join(mochaDir, "apps", app, "current", "install.json")
		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			return fmt.Errorf("install info file does not exist for %s", app)
		}

		infoFile, err := os.ReadFile(infoPath)
		if err != nil {
			return fmt.Errorf("cannot read install info file for %s: %w", app, err)
		}

		installInfo := pkg.InstallInfo{}

		if err := json.Unmarshal(infoFile, &installInfo); err != nil {
			return fmt.Errorf("cannot parse install info file for %s: %w", app, err)
		}

		manifestRef, err := manifest.PopulateRef(manifest.Ref{Name: app, Bucket: installInfo.Bucket}, mochaDir)
		if err != nil {
			return fmt.Errorf("cannot fetch app info for %s: %w", app, err)
		}

		if manifest.CompareVersions(installInfo.Version, manifestRef.Version) != 1 {
			continue
		}

		output.LogOutput("upgrading " + app)

		oldVersionDir := filepath.Join(mochaDir, "apps", app, installInfo.Version)
		if err := pkg.UnlinkApp(manifestRef, downloadArch, mochaDir, oldVersionDir); err != nil {
			return fmt.Errorf("failed to unlink old version of %s: %w", app, err)
		}

		if err := pkg.InstallApp(manifestRef, downloadArch, false, mochaDir); err != nil {
			return fmt.Errorf("failed to install %s: %w", app, err)
		}

		output.LogOutput(fmt.Sprintf("Upgraded %s", app))
	}

	return nil
}
