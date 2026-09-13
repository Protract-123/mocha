package shim

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/output"
)

type Info struct {
	Name             string
	Target           string
	Args             string
	WorkingDirectory string
	Elevate          bool
	EnvVars          map[string]string
}

func CreateShim(info Info, mochaDir string) error {
	if info.Name == "" {
		return fmt.Errorf("missing shim name")
	} else if info.Target == "" || !filepath.IsAbs(info.Target) {
		return fmt.Errorf("missing shim target")
	}

	shimDir := filepath.Join(mochaDir, "shims")
	if err := os.MkdirAll(shimDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create shim directory: %w", err)
	}

	fileExtension := strings.ToLower(filepath.Ext(info.Target))
	switch {
	case fileExtension == ".exe" || fileExtension == ".com":
		break
	case fileExtension == ".bat" || fileExtension == ".cmd":
		cmdPath, err := exec.LookPath("cmd.exe")
		if err != nil {
			return fmt.Errorf("failed to find cmd executable: %w", err)
		}
		info.Args = "/C " + `"` + info.Target + `"` + " " + info.Args
		info.Target = cmdPath
	case fileExtension == ".ps1":
		// TODO: make powershell interpreter name/path customizable
		powershellPath, err := exec.LookPath("powershell.exe")
		if err != nil {
			return fmt.Errorf("failed to find cmd executable: %w", err)
		}
		info.Args = "-File " + `"` + info.Target + `"` + " " + info.Args
		info.Target = powershellPath
	case fileExtension == ".jar":
		// TODO: make java runtime name/path customizable
		powershellPath, err := exec.LookPath("java.exe")
		if err != nil {
			return fmt.Errorf("failed to find cmd executable: %w", err)
		}
		info.Args = "-jar " + `"` + info.Target + `"` + " " + info.Args
		info.Target = powershellPath

	case fileExtension == ".py":
		// TODO: make python interpreter name/path customizable
		pythonPath, err := exec.LookPath("python3.14.exe")
		if err != nil {
			return fmt.Errorf("failed to find python executable: %w", err)
		}
		info.Args = `"` + info.Target + `"` + " " + info.Args
		info.Target = pythonPath
	default:
		output.LogWarning("unknown extension %q, skipping shim", fileExtension)
		return nil
	}

	if err := copyShimBinary(info, true, mochaDir); err != nil {
		return fmt.Errorf("failed to create exe shim: %w", err)
	}

	if err := createShimSidecar(info, shimDir); err != nil {
		return fmt.Errorf("failed to create .shim sidecar file: %w", err)
	}

	return nil
}

func DeleteShim(name string, mochaDir string) error {
	shimsDir := filepath.Join(mochaDir, "shims")

	shims, err := os.ReadDir(shimsDir)
	if err != nil {
		return fmt.Errorf("failed to read shims directory: %w", err)
	}

	deletion := false

	for _, shim := range shims {
		if shim.IsDir() {
			continue
		}

		shimName := strings.TrimSuffix(shim.Name(), filepath.Ext(shim.Name()))

		if shimName == name {
			if err := os.Remove(filepath.Join(shimsDir, shim.Name())); err != nil {
				return fmt.Errorf("failed to remove shim file %s: %w", shim.Name(), err)
			}
			deletion = true
		}
	}

	if !deletion {
		output.LogWarning("no shim found for %q", name)
	}

	return nil
}

func GetAllShims(mochaDir string) ([]Info, error) {
	shimsDir := filepath.Join(mochaDir, "shims")

	shims, err := os.ReadDir(shimsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read shims directory: %w", err)
	}

	var shimInfo []Info

	for _, shim := range shims {
		if shim.IsDir() {
			continue
		}

		if filepath.Ext(shim.Name()) == ".shim" {
			info, err := parseShimSidecar(filepath.Join(shimsDir, shim.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", shim.Name(), err)
			}

			shimInfo = append(shimInfo, info)
		}
	}

	return shimInfo, nil
}
