package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Protract-123/mocha/config"
)

type Command struct{}

func (cmd *Command) Run() error {
	configPath, err := config.Location(config.Current().MochaDirectory)
	if errors.Is(err, config.ErrNotFound) {
		if err := config.WriteDefault(configPath); err != nil {
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
		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		err = editorCmd.Run()
	} else {
		err = exec.Command("cmd.exe", "/c", "start", configPath).Run()
	}

	if err != nil {
		return fmt.Errorf("failed to open config file in editor: %w", err)
	}

	return nil
}
