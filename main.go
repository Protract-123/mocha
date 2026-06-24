package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Protract-123/mocha/commands"
	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/output"
	"github.com/alexflint/go-arg"
)

var args struct {
	MochaDir string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha"`

	Bucket   *commands.BucketCommand   `arg:"subcommand:bucket"`
	Cache    *commands.CacheCommand    `arg:"subcommand:cache"`
	Cat      *commands.CatCommand      `arg:"subcommand:cat"`
	Config   *commands.ConfigCommand   `arg:"subcommand:config"`
	Download *commands.DownloadCommand `arg:"subcommand:download"`
	Search   *commands.SearchCommand   `arg:"subcommand:search"`
	Shim     *commands.ShimCommand     `arg:"subcommand:shim"`
	Update   *commands.UpdateCommand   `arg:"subcommand:update"`
}

type Config struct {
	CatConfig commands.CatConfig `toml:"cat"`
}

func main() {
	if err := run(); err != nil {
		output.LogError(err)
		os.Exit(1)
	}
}

func run() error {
	arg.MustParse(&args)

	mochaDir := os.ExpandEnv(args.MochaDir)
	if !filepath.IsAbs(mochaDir) {
		return fmt.Errorf("mocha dir %q is not an absolute path", mochaDir)
	}

	if err := os.MkdirAll(mochaDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create mocha dir: %w", err)
	}

	appConfig, err := getConfig(mochaDir)
	if err != nil {
		return err
	}

	switch {
	case args.Bucket != nil:
		return args.Bucket.Run(mochaDir)
	case args.Cache != nil:
		return args.Cache.Run(mochaDir)
	case args.Cat != nil:
		return args.Cat.Run(mochaDir, appConfig.CatConfig)
	case args.Config != nil:
		return args.Config.Run(mochaDir)
	case args.Download != nil:
		return args.Download.Run(mochaDir)
	case args.Search != nil:
		return args.Search.Run(mochaDir)
	case args.Shim != nil:
		return args.Shim.Run(mochaDir)
	case args.Update != nil:
		return args.Update.Run(mochaDir)
	}
	return nil
}

func getConfig(mochaDir string) (*Config, error) {
	appConfig := &Config{}

	configPath, err := config.GetConfigPath(mochaDir)
	if errors.Is(err, config.ConfigNotFound) {
		output.LogWarning(fmt.Sprintf("failed to find mocha.toml, using defaults"))
		return appConfig, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err = toml.DecodeFile(configPath, appConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return appConfig, nil
}
