package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
	"github.com/Protract-123/mocha/output"
	"github.com/Protract-123/mocha/pkg"
	"github.com/alexflint/go-arg"
)

type UpgradeCommand struct {
	Apps       []string `arg:"positional" help:"apps to upgrade (e.g. git, 7zip)"`
	All        bool     `arg:"-a,--all" help:"upgrade all apps"`
	Force      bool     `arg:"-f,--force" help:"ignore cache hits"`
	SkipVerify bool     `arg:"-s,--skip-verify" help:"skip hash check"`
}

func (cmd *UpgradeCommand) Run() error {
	if len(cmd.Apps) == 0 && !cmd.All {
		return arg.ErrHelp
	}

	downloadArch, err := pkg.DownloadArch()
	if err != nil {
		return fmt.Errorf("failed to get system architecture: %w", err)
	}

	mochaDir := config.Current().MochaDirectory

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

		info, err := manifest.PopulateInfo(manifest.Info{Name: app, Bucket: installInfo.Bucket}, mochaDir)
		if err != nil {
			return fmt.Errorf("cannot fetch app info for %s: %w", app, err)
		}

		if manifest.CompareVersions(installInfo.Version, info.Version) != 1 {
			continue
		}

		manifestJson, err := manifest.GetJson(info.ManifestPath)
		if err != nil {
			return fmt.Errorf("failed to get manifest JSON: %w", err)
		}

		output.LogInfo("upgrading " + app)

		target := pkg.Package{
			Info: info,
			Json: manifestJson,
			Arch: downloadArch,
		}

		options := pkg.DownloadOptions{
			Force:      cmd.Force,
			SkipVerify: cmd.SkipVerify,
		}

		downloadResults, err := pkg.Download(target, mochaDir, options)
		if err != nil {
			return fmt.Errorf("failed to download manifest files: %w", err)
		}

		if err := pkg.Install(target, downloadResults, mochaDir); err != nil {
			return fmt.Errorf("failed to install %s: %w", app, err)
		}

		if err := pkg.Unlink(app, mochaDir); err != nil {
			return fmt.Errorf("failed to unlink old version of %s: %w", app, err)
		}

		if err := pkg.Link(target.Info, mochaDir); err != nil {
			return fmt.Errorf("failed to link app: %w", err)
		}

		output.LogSuccess("successfully upgraded %q", app)
	}

	return nil
}
