package commands

import (
	_ "embed"
	"errors"
	"os"
	"os/exec"

	"github.com/Protract-123/mocha/config"
)

type ConfigCommand struct{}

func (cmd ConfigCommand) Run(mochaDir string) error {
	configPath, err := config.GetConfigPath(mochaDir)
	if errors.Is(err, config.ConfigNotFound) {
		err := config.WriteDefaultConfig(configPath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
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
		return cmd.Run()
	}

	return exec.Command("cmd", "/c", "start", configPath).Run()
}
