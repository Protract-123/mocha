package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/manifest"
)

type CatCommand struct {
	Name   string `arg:"positional,required"`
	Bucket string `arg:"-b, --bucket"`
}

type CatConfig struct {
	IncludeDeprecated bool   `toml:"include-deprecated"`
	Command           string `toml:"command"`
}

func (cmd CatCommand) Run(mochaDir string, config CatConfig) error {
	writeManifestData := func(path string) error {
		if config.Command == "" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(data)
			if err != nil {
				return err
			}

			return nil
		}

		if !strings.Contains(config.Command, "[path]") {
			return fmt.Errorf("command %s must contain [path] to replace", config.Command)
		}

		commandStr := strings.Replace(config.Command, "[path]", path, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		return command.Run()
	}

	if cmd.Bucket != "" {
		path, err := manifest.GetManifestPath(cmd.Bucket, cmd.Name, mochaDir)
		if os.IsNotExist(err) {
			return fmt.Errorf("app %s not found in bucket %s", cmd.Name, cmd.Bucket)
		} else if err != nil {
			return err
		}

		return writeManifestData(path)
	}

	bucketsDir := filepath.Join(mochaDir, "buckets")
	buckets, err := os.ReadDir(bucketsDir)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if !bucket.IsDir() {
			continue
		}

		path, err := manifest.GetManifestPath(bucket.Name(), cmd.Name, mochaDir)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}

		return writeManifestData(path)
	}

	return fmt.Errorf("app %s not found in buckets", cmd.Name)
}
