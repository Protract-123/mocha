package config

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestWriteDefault(tester *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid path: writable directory",
			args:    args{path: filepath.Join(tester.TempDir(), "mocha.toml")},
			wantErr: false,
		},
		{
			name:    "invalid path: parent directory does not exist",
			args:    args{path: filepath.Join(tester.TempDir(), "missing", "mocha.toml")},
			wantErr: true,
		},
	}
	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			err := WriteDefault(test.args.path)
			if (err != nil) != test.wantErr {
				t.Errorf("WriteDefault() error = %v, wantErr %v", err, test.wantErr)
			}

			if test.wantErr {
				return
			}

			written, err := os.ReadFile(test.args.path)
			if err != nil {
				t.Fatalf("failed to read back written config: %v", err)
			}

			if !bytes.Equal(defaultConfigToml, written) {
				t.Errorf("WriteDefault() wrote %q, want %q", written, defaultConfigToml)
			}
		})
	}
}

func TestDefaultMatchesFile(tester *testing.T) {
	var configFromFile MochaConfiguration
	if _, err := toml.DecodeFile("default_config.toml", &configFromFile); err != nil {
		tester.Fatalf("failed to decode default_config.toml: %v", err)
	}

	if !reflect.DeepEqual(configFromFile, Default()) {
		tester.Errorf("default_config.toml is out of sync with DefaultConfig()")
	}
}
