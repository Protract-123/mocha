package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//goland:noinspection GoNameStartsWithPackageName
var ConfigNotFound = errors.New("mocha.toml not found")

//go:embed default_config.toml
var defaultConfig []byte

func GetConfigPath(mochaDir string) (string, error) {
	configDirs := []string{
		mochaDir,
		filepath.Join(os.ExpandEnv("$APPDATA"), "mocha"),
		filepath.Join(os.ExpandEnv("$XDG_CONFIG_HOME"), "mocha"),
		filepath.Join(os.ExpandEnv("$USERPROFILE"), ".config", "mocha"),
	}

	for _, dir := range configDirs {
		if !filepath.IsAbs(dir) {
			continue
		}

		configPath := filepath.Join(dir, "mocha.toml")

		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("failed to check if config exists: %w", err)
		}

		return configPath, nil
	}

	return filepath.Join(mochaDir, "mocha.toml"), ConfigNotFound
}

func WriteDefaultConfig(configPath string) error {
	err := os.WriteFile(configPath, defaultConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
