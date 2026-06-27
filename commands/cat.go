package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Protract-123/mocha/manifest"
)

type CatCommand struct {
	ManifestReferences []string `arg:"positional,required"`
}

type CatConfig struct {
	IncludeDeprecated bool   `toml:"include-deprecated"`
	Command           string `toml:"command"`
}

func (cmd CatCommand) Run(mochaDir string, config CatConfig) error {
	if cmd.ManifestReferences == nil {
		return errors.New("at least one manifest reference is required")
	}

	for _, refString := range cmd.ManifestReferences {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, mochaDir)
		if err != nil {
			return fmt.Errorf("failed to get %q manifest details: %w", refString, err)
		}

		if config.Command == "" {
			data, err := os.ReadFile(manifestRef.ManifestPath)
			if err != nil {
				return fmt.Errorf("failed to read manifest %q: %w", manifestRef.ManifestPath, err)
			}

			_, err = os.Stdout.Write(data)
			if err != nil {
				return fmt.Errorf("failed to display manifest %q: %w", manifestRef.ManifestPath, err)
			}

			continue
		}

		if !strings.Contains(config.Command, "[path]") {
			return fmt.Errorf("command %s must contain [path] to replace", config.Command)
		}

		commandStr := strings.Replace(config.Command, "[path]", manifestRef.ManifestPath, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()

		if err != nil {
			return fmt.Errorf("failed to display manifest %q: %w", manifestRef.ManifestPath, err)
		}
	}

	return nil
}
