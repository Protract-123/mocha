package config

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed default_config.toml
var defaultConfigToml []byte

func Default() MochaConfiguration {
	return MochaConfiguration{
		Cat: CatConfig{
			IncludeDeprecated: false,
			Command:           "",
		},
		Colors: ColorConfig{
			SuccessColor: "magenta",
			ErrorColor:   "red",
			WarningColor: "yellow",
			InfoColor:    "blue",
		},
		MochaDirectory: "",
	}
}

func WriteDefault(path string) error {
	if err := os.WriteFile(path, defaultConfigToml, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
