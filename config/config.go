package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/Protract-123/mocha/output"
)

type MochaConfiguration struct {
	CatConfiguration CatConfig `toml:"cat"`
}

type CatConfig struct {
	IncludeDeprecated bool   `toml:"include-deprecated"`
	Command           string `toml:"command"`
}

func GetConfig(mochaDir string) (*MochaConfiguration, error) {
	appConfig := &MochaConfiguration{}

	configPath, err := GetConfigPath(mochaDir)
	if errors.Is(err, ConfigNotFound) {
		output.LogWarning(fmt.Sprintf("failed to find mocha.toml, using defaults"))
		return appConfig, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err = toml.DecodeFile(configPath, appConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return appConfig, nil
}
