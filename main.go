package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Protract-123/mocha/commands"
	"github.com/Protract-123/mocha/config"
	"github.com/alexflint/go-arg"
)

var args struct {
	MochaDir string `arg:"--,env:MOCHA_DIR" default:"$USERPROFILE/mocha"`

	Bucket   *commands.BucketCommand   `arg:"subcommand:bucket"`
	Cat      *commands.CatCommand      `arg:"subcommand:cat"`
	Config   *commands.ConfigCommand   `arg:"subcommand:config"`
	Download *commands.DownloadCommand `arg:"subcommand:download"`
}

type Config struct {
	CatConfig commands.CatConfig `toml:"cat"`
}

func main() {
	arg.MustParse(&args)

	mochaDir, err := filepath.Abs(os.ExpandEnv(args.MochaDir))
	if err != nil {
		log.Fatalf("failed to resolve mocha dir: %v", err)
	}

	configPath, err := config.GetConfigPath(mochaDir)
	if err != nil {
		log.Fatalf("failed to get config path: %v", err)
	}

	appConfig := &Config{}

	_, err = toml.DecodeFile(configPath, appConfig)
	if err != nil {
		return
	}

	if args.Bucket != nil {
		err := args.Bucket.Run(mochaDir)
		if err != nil {
			println(err.Error())
		}
	}
	if args.Cat != nil {
		err := args.Cat.Run(mochaDir, appConfig.CatConfig)
		if err != nil {
			println(err.Error())
		}
	}
	if args.Config != nil {
		err := args.Config.Run(mochaDir)
		if err != nil {
			println(err.Error())
		}
	}
	if args.Download != nil {
		err := args.Download.Run(mochaDir)
		if err != nil {
			println(err.Error())
		}
	}
}
