package config

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed default_config.toml
var defaultConfigToml []byte

var defaultConfig = MochaConfiguration{
	Cat: CatConfig{
		IncludeDeprecated: false,
		Command:           "",
	},
	Colors: ColorConfig{
		AccentColor: "purple",
		ErrorColor:  "red",
		WarnColor:   "yellow",
		InfoColor:   "blue",
	},
	MochaDirectory: "",
}

func WriteDefault(path string) error {
	if err := os.WriteFile(path, defaultConfigToml, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
