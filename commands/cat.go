package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
)

type CatCommand struct {
	Manifests []string `arg:"positional,required" help:"manifests to show (e.g. git, zed)"`
}

func (cmd *CatCommand) Run() error {
	catConfig := config.Current().Cat

	for _, refString := range cmd.Manifests {
		manifestRef, err := manifest.ParseRefString(refString)
		if err != nil {
			return fmt.Errorf("failed to parse manifest ref %q: %w", refString, err)
		}

		manifestRef, err = manifest.PopulateRef(manifestRef, config.Current().MochaDirectory)
		if err != nil {
			return fmt.Errorf("failed to get manifest details for %q: %w", refString, err)
		}

		if catConfig.Command == "" {
			data, err := os.ReadFile(manifestRef.ManifestPath)
			if err != nil {
				return fmt.Errorf("failed to read manifest %q: %w", manifestRef.ManifestPath, err)
			}

			if _, err := os.Stdout.Write(data); err != nil {
				return fmt.Errorf("failed to display manifest %q: %w", manifestRef.ManifestPath, err)
			}

			continue
		}

		if !strings.Contains(catConfig.Command, "[path]") {
			return fmt.Errorf("command %q must contain [path] to replace", catConfig.Command)
		}

		commandStr := strings.Replace(catConfig.Command, "[path]", manifestRef.ManifestPath, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to display manifest %q: %w", manifestRef.ManifestPath, err)
		}
	}

	return nil
}
