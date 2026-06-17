package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Protract-123/mocha/app"
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

	for _, entry := range cmd.AppReferences {
		appRef, err := app.ParseAppString(entry)
		if err != nil {
			return err
		}

		appRef, err = app.PopulateAppRef(appRef, mochaDir)
		if err != nil {
			return err
		}

		path, err := app.GetManifestPath(appRef.Bucket, appRef.Name, mochaDir)
		if os.IsNotExist(err) {
			return fmt.Errorf("app %s not found in bucket %s", appRef.Name, appRef.Bucket)
		} else if err != nil {
			return err
		}

		if config.Command == "" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(data)
			if err != nil {
				return err
			}

			continue
		}

		if !strings.Contains(config.Command, "[path]") {
			return fmt.Errorf("command %s must contain [path] to replace", config.Command)
		}

		commandStr := strings.Replace(config.Command, "[path]", path, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()

		if err != nil {
			return err
		}
	}

	return nil
}
