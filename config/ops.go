package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var ErrNotFound = errors.New("mocha.toml not found")

var (
	currentConfig  MochaConfiguration
	loadConfigOnce sync.Once
)

func Init(config MochaConfiguration) {
	loadConfigOnce.Do(func() {
		currentConfig = config
	})
}

func Current() MochaConfiguration {
	return currentConfig
}

func Load(mochaDir string) (MochaConfiguration, error) {
	config := MochaConfiguration{}
	configPath, err := Location(mochaDir)

	if errors.Is(err, ErrNotFound) {
		return Default(), ErrNotFound
	} else if err != nil {
		return MochaConfiguration{}, fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err = toml.DecodeFile(configPath, &config); err != nil {
		return MochaConfiguration{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

func Location(mochaDir string) (string, error) {
	configLocations := []string{
		filepath.Join(mochaDir, "mocha.toml"),
		filepath.Join(os.ExpandEnv("$APPDATA"), "mocha", "mocha.toml"),
		filepath.Join(os.ExpandEnv("$XDG_CONFIG_HOME"), "mocha", "mocha.toml"),
		filepath.Join(os.ExpandEnv("$USERPROFILE"), ".config", "mocha", "mocha.toml"),
	}

	for _, path := range configLocations {
		if !filepath.IsAbs(path) {
			continue
		}

		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("failed to check if config exists: %w", err)
		}

		return path, nil
	}

	return filepath.Join(mochaDir, "mocha.toml"), ErrNotFound
}
