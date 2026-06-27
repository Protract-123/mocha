package commands

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Protract-123/mocha/config"
)

type ConfigCommand struct{}

func (cmd ConfigCommand) Run(mochaDir string) error {
	configPath, err := config.GetConfigPath(mochaDir)
	if errors.Is(err, config.ErrConfigNotFound) {
		err := config.WriteDefaultConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to write default config: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if config file exists: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}

	if editor != "" {
		cmd := exec.Command(editor, configPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	} else {
		err = exec.Command("cmd", "/c", "start", configPath).Run()
	}

	if err != nil {
		return fmt.Errorf("failed to open config file in editor: %w", err)
	}

	return nil
}
