package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/commands"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

type arguments struct {
	MochaDirectory string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha" help:"directory where mocha stores buckets, apps, shims, and cache"`

	BucketCommand    *commands.BucketCommand    `arg:"subcommand:bucket" help:"manage mocha buckets"`
	CacheCommand     *commands.CacheCommand     `arg:"subcommand:cache" help:"manage cache items"`
	CatCommand       *commands.CatCommand       `arg:"subcommand:cat" help:"show an app's manifest"`
	ConfigCommand    *commands.ConfigCommand    `arg:"subcommand:config" help:"edit the config file"`
	DownloadCommand  *commands.DownloadCommand  `arg:"subcommand:download" help:"download and verify an app's files"`
	InstallCommand   *commands.InstallCommand   `arg:"subcommand:install" help:"install apps"`
	ListCommand      *commands.ListCommand      `arg:"subcommand:list" help:"list installed apps"`
	SearchCommand    *commands.SearchCommand    `arg:"subcommand:search" help:"search for an app in buckets"`
	ShimCommand      *commands.ShimCommand      `arg:"subcommand:shim" help:"manage mocha shims"`
	UninstallCommand *commands.UninstallCommand `arg:"subcommand:uninstall" help:"uninstall apps"`
	UpdateCommand    *commands.UpdateCommand    `arg:"subcommand:update" help:"update mocha buckets"`
	UpgradeCommand   *commands.UpgradeCommand   `arg:"subcommand:upgrade" help:"upgrade mocha apps"`
}

func (arguments) Version() string {
	return "mocha v0.0.1"
}

func main() {
	if err := run(); err != nil {
		output.LogError(err)
		os.Exit(1)
	}
}

func run() error {
	var args arguments
	parser := arg.MustParse(&args)

	mochaDirectory := os.ExpandEnv(args.MochaDirectory)
	if !filepath.IsAbs(mochaDirectory) {
		return fmt.Errorf("mocha directory %q is not an absolute path", mochaDirectory)
	}

	if err := os.MkdirAll(mochaDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create mocha directory: %w", err)
	}

	configuration, err := config.Load(mochaDirectory)
	if err != nil && !errors.Is(err, config.ErrNotFound) {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	configuration.MochaDirectory = mochaDirectory
	config.Init(configuration)

	if errors.Is(err, config.ErrNotFound) {
		output.LogWarning("failed to find mocha.toml, using defaults")
	}

	var executionError error

	switch {
	case args.BucketCommand != nil:
		executionError = args.BucketCommand.Run(mochaDirectory)
	case args.CacheCommand != nil:
		executionError = args.CacheCommand.Run(mochaDirectory)
	case args.CatCommand != nil:
		executionError = args.CatCommand.Run(mochaDirectory, configuration.Cat)
	case args.ConfigCommand != nil:
		executionError = args.ConfigCommand.Run(mochaDirectory)
	case args.DownloadCommand != nil:
		executionError = args.DownloadCommand.Run(mochaDirectory)
	case args.InstallCommand != nil:
		executionError = args.InstallCommand.Run(mochaDirectory)
	case args.ListCommand != nil:
		executionError = args.ListCommand.Run(mochaDirectory)
	case args.SearchCommand != nil:
		executionError = args.SearchCommand.Run(mochaDirectory)
	case args.ShimCommand != nil:
		executionError = args.ShimCommand.Run(mochaDirectory)
	case args.UninstallCommand != nil:
		executionError = args.UninstallCommand.Run(mochaDirectory)
	case args.UpdateCommand != nil:
		executionError = args.UpdateCommand.Run(mochaDirectory)
	case args.UpgradeCommand != nil:
		executionError = args.UpgradeCommand.Run(mochaDirectory)
	default:
		executionError = arg.ErrHelp
	}

	if errors.Is(executionError, arg.ErrHelp) {
		parser.WriteHelp(os.Stdout)
		return nil
	}
	return executionError
}
