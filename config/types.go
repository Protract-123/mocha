package config

type MochaConfiguration struct {
	Cat    CatConfig   `toml:"cat"`
	Colors ColorConfig `toml:"colors"`

	MochaDirectory string `toml:"-"`
}

type CatConfig struct {
	IncludeDeprecated bool   `toml:"include-deprecated"`
	Command           string `toml:"command"`
}

type ColorConfig struct {
	AccentColor  string `toml:"accent"`
	ErrorColor   string `toml:"error"`
	WarningColor string `toml:"warning"`
	InfoColor    string `toml:"info"`
}
