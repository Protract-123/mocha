package cat

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Protract-123/mocha/config"
	"github.com/Protract-123/mocha/manifest"
)

type Command struct {
	Manifests []string `arg:"positional,required" help:"manifests to show (e.g. git, zed)"`
}

func (cmd *Command) Run() error {
	catConfig := config.Current().Cat

	for _, spec := range cmd.Manifests {
		info, err := manifest.ParseSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse manifest spec %q: %w", spec, err)
		}

		info, err = manifest.PopulateInfo(info, config.Current().MochaDirectory)
		if err != nil {
			return fmt.Errorf("failed to get manifest details for %q: %w", spec, err)
		}

		if catConfig.Command == "" {
			data, err := os.ReadFile(info.ManifestPath)
			if err != nil {
				return fmt.Errorf("failed to read manifest %q: %w", info.ManifestPath, err)
			}

			if _, err := os.Stdout.Write(data); err != nil {
				return fmt.Errorf("failed to display manifest %q: %w", info.ManifestPath, err)
			}

			continue
		}

		if !strings.Contains(catConfig.Command, "[path]") {
			return fmt.Errorf("command %q must contain [path] to replace", catConfig.Command)
		}

		commandStr := strings.Replace(catConfig.Command, "[path]", info.ManifestPath, 1)

		command := exec.Command("cmd.exe", "/C", commandStr)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to display manifest %q: %w", info.ManifestPath, err)
		}
	}

	return nil
}
