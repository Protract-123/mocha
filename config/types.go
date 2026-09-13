package config

type MochaConfiguration struct {
	Cat    CatConfig   `toml:"cat"`
	Colors ColorConfig `toml:"colors"`
	Shim   ShimConfig  `toml:"shim"`

	MochaDirectory string `toml:"-"`
}

type CatConfig struct {
	Command string `toml:"command"`
}

type ColorConfig struct {
	SuccessColor string `toml:"success"`
	ErrorColor   string `toml:"error"`
	WarningColor string `toml:"warning"`
	InfoColor    string `toml:"info"`
}

type ShimConfig struct {
	JavaRunner       string `toml:"java-runner"`
	PowerShellRunner string `toml:"powershell-runner"`
	PythonRunner     string `toml:"python-runner"`
}
