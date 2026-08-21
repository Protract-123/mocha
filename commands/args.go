package commands

import (
	"github.com/Protract-123/mocha/commands/bucket"
	"github.com/Protract-123/mocha/commands/cache"
	"github.com/Protract-123/mocha/commands/cat"
	"github.com/Protract-123/mocha/commands/config"
	"github.com/Protract-123/mocha/commands/download"
	"github.com/Protract-123/mocha/commands/install"
	"github.com/Protract-123/mocha/commands/list"
	"github.com/Protract-123/mocha/commands/outdated"
	"github.com/Protract-123/mocha/commands/search"
	"github.com/Protract-123/mocha/commands/shim"
	"github.com/Protract-123/mocha/commands/uninstall"
	"github.com/Protract-123/mocha/commands/update"
	"github.com/Protract-123/mocha/commands/upgrade"
)

type Arguments struct {
	MochaDirectory string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha" help:"directory where mocha stores buckets, apps, shims, and cache"`

	BucketCommand    *bucket.Command    `arg:"subcommand:bucket" help:"manage mocha buckets"`
	CacheCommand     *cache.Command     `arg:"subcommand:cache" help:"manage cache items"`
	CatCommand       *cat.Command       `arg:"subcommand:cat" help:"show an app's manifest"`
	ConfigCommand    *config.Command    `arg:"subcommand:config" help:"edit the config file"`
	DownloadCommand  *download.Command  `arg:"subcommand:download" help:"download and verify an app's files"`
	InstallCommand   *install.Command   `arg:"subcommand:install" help:"install apps"`
	ListCommand      *list.Command      `arg:"subcommand:list" help:"list installed apps"`
	OutdatedCommand  *outdated.Command  `arg:"subcommand:outdated" help:"list outdated apps"`
	SearchCommand    *search.Command    `arg:"subcommand:search" help:"search for an app in buckets"`
	ShimCommand      *shim.Command      `arg:"subcommand:shim" help:"manage mocha shims"`
	UninstallCommand *uninstall.Command `arg:"subcommand:uninstall" help:"uninstall apps"`
	UpdateCommand    *update.Command    `arg:"subcommand:update" help:"update mocha buckets"`
	UpgradeCommand   *upgrade.Command   `arg:"subcommand:upgrade" help:"upgrade mocha apps"`
}

func (Arguments) Version() string {
	return "mocha v0.0.1"
}
