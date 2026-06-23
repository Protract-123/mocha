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
			return fmt.Errorf("failed to parse app ref %q: %w", entry, err)
		}

		appRef, err = app.PopulateAppRef(appRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", entry, err)
		}

		if config.Command == "" {
			data, err := os.ReadFile(appRef.ManifestPath)
			if err != nil {
				return fmt.Errorf("failed to read manifest %q: %w", appRef.ManifestPath, err)
			}

			_, err = os.Stdout.Write(data)
			if err != nil {
				return fmt.Errorf("failed to display manifest %q: %w", appRef.ManifestPath, err)
			}

			continue
		}

		if !strings.Contains(config.Command, "[path]") {
			return fmt.Errorf("command %s must contain [path] to replace", config.Command)
		}

		commandStr := strings.Replace(config.Command, "[path]", appRef.ManifestPath, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()

		if err != nil {
			return fmt.Errorf("failed to display manifest %q: %w", appRef.ManifestPath, err)
		}
	}

	return nil
}
