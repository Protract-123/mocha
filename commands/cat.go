package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/bucket"
	"github.com/Protract-123/mocha/manifest"
)

type CatCommand struct {
	AppReferences []string `arg:"positional,required"`
}

type CatConfig struct {
	IncludeDeprecated bool   `toml:"include-deprecated"`
	Command           string `toml:"command"`
}

func (cmd CatCommand) Run(mochaDir string, config CatConfig) error {
	if cmd.AppReferences == nil {
		return errors.New("at least one app reference is required")
	}

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

	bucketsDir := filepath.Join(mochaDir, "buckets")
	buckets, err := os.ReadDir(bucketsDir)
	if err != nil {
		return err
	}

	for _, entry := range cmd.AppReferences {
		appRef, err := bucket.ParseAppString(entry)
		if err != nil {
			return err
		}

		if appRef.Bucket != "" {
			path, err := manifest.GetManifestPath(appRef.Bucket, appRef.Name, mochaDir)
			if os.IsNotExist(err) {
				return fmt.Errorf("app %s not found in bucket %s", appRef.Name, appRef.Bucket)
			} else if err != nil {
				return err
			}

			err = writeManifestData(path)
			if err != nil {
				return err
			}

			continue
		}

		found := false

		for _, dirEntry := range buckets {
			if !dirEntry.IsDir() {
				continue
			}

			path, err := manifest.GetManifestPath(dirEntry.Name(), appRef.Name, mochaDir)
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return err
			}

			err = writeManifestData(path)
			if err != nil {
				return err
			}

			found = true
			break
		}

		if found {
			continue
		}

		return fmt.Errorf("app %s not found in buckets", appRef.Name)
	}

	return nil
}
