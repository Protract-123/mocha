package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protract-123/mocha/commands"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

var args struct {
	MochaDirectory string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha"`

	BucketCommand    *commands.BucketCommand    `arg:"subcommand:bucket"`
	CacheCommand     *commands.CacheCommand     `arg:"subcommand:cache"`
	CatCommand       *commands.CatCommand       `arg:"subcommand:cat"`
	ConfigCommand    *commands.ConfigCommand    `arg:"subcommand:config"`
	DownloadCommand  *commands.DownloadCommand  `arg:"subcommand:download"`
	InstallCommand   *commands.InstallCommand   `arg:"subcommand:install"`
	SearchCommand    *commands.SearchCommand    `arg:"subcommand:search"`
	ShimCommand      *commands.ShimCommand      `arg:"subcommand:shim"`
	UninstallCommand *commands.UninstallCommand `arg:"subcommand:uninstall"`
	UpdateCommand    *commands.UpdateCommand    `arg:"subcommand:update"`
}

func main() {
	if err := run(); err != nil {
		output.LogError(err)
		os.Exit(1)
	}
}

func run() error {
	arg.MustParse(&args)

	mochaDirectory := os.ExpandEnv(args.MochaDirectory)
	if !filepath.IsAbs(mochaDirectory) {
		return fmt.Errorf("mocha directory %q is not an absolute path", mochaDirectory)
	}

	if err := os.MkdirAll(mochaDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create mocha directory: %w", err)
	}

	configuration, err := config.GetConfig(mochaDirectory)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	switch {
	case args.BucketCommand != nil:
		return args.BucketCommand.Run(mochaDirectory)
	case args.CacheCommand != nil:
		return args.CacheCommand.Run(mochaDirectory)
	case args.CatCommand != nil:
		return args.CatCommand.Run(mochaDirectory, configuration.CatConfiguration)
	case args.ConfigCommand != nil:
		return args.ConfigCommand.Run(mochaDirectory)
	case args.DownloadCommand != nil:
		return args.DownloadCommand.Run(mochaDirectory)
	case args.InstallCommand != nil:
		return args.InstallCommand.Run(mochaDirectory)
	case args.SearchCommand != nil:
		return args.SearchCommand.Run(mochaDirectory)
	case args.ShimCommand != nil:
		return args.ShimCommand.Run(mochaDirectory)
	case args.UninstallCommand != nil:
		return args.UninstallCommand.Run(mochaDirectory)
	case args.UpdateCommand != nil:
		return args.UpdateCommand.Run(mochaDirectory)
	}
	return nil
}
