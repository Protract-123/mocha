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

func main() {
	if err := run(); err != nil {
		output.LogError(err)
		os.Exit(1)
	}
}

func run() error {
	var args commands.Arguments
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
	output.InitTheme()

	if errors.Is(err, config.ErrNotFound) {
		output.LogWarning("failed to find mocha.toml, using defaults")
	}

	var executionError error

	switch {
	case args.BucketCommand != nil:
		executionError = args.BucketCommand.Run()
	case args.CacheCommand != nil:
		executionError = args.CacheCommand.Run()
	case args.CatCommand != nil:
		executionError = args.CatCommand.Run()
	case args.ConfigCommand != nil:
		executionError = args.ConfigCommand.Run()
	case args.DownloadCommand != nil:
		executionError = args.DownloadCommand.Run()
	case args.InstallCommand != nil:
		executionError = args.InstallCommand.Run()
	case args.ListCommand != nil:
		executionError = args.ListCommand.Run()
	case args.OutdatedCommand != nil:
		executionError = args.OutdatedCommand.Run()
	case args.SearchCommand != nil:
		executionError = args.SearchCommand.Run()
	case args.ShimCommand != nil:
		executionError = args.ShimCommand.Run()
	case args.UninstallCommand != nil:
		executionError = args.UninstallCommand.Run()
	case args.UpdateCommand != nil:
		executionError = args.UpdateCommand.Run()
	case args.UpgradeCommand != nil:
		executionError = args.UpgradeCommand.Run()
	default:
		executionError = arg.ErrHelp
	}

	if errors.Is(executionError, arg.ErrHelp) {
		parser.WriteHelp(os.Stdout)
		return nil
	}
	return executionError
}
