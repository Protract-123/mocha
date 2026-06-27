package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed default_config.toml
var defaultConfig []byte

var ErrConfigNotFound = errors.New("mocha.toml not found")

func GetConfigPath(mochaDir string) (string, error) {
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

	return filepath.Join(mochaDir, "mocha.toml"), ErrConfigNotFound
}

func WriteDefaultConfig(configPath string) error {
	if err := os.WriteFile(configPath, defaultConfig, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
