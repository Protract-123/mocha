package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createShimSidecar(info Info, shimDir string) error {
	if info.Name == "" {
		return fmt.Errorf("missing shim name")
	} else if info.Target == "" {
		return fmt.Errorf("missing shim target")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("path = %s\n", info.Target))

	if info.Args != "" {
		sb.WriteString(fmt.Sprintf("args = %s\n", info.Args))
	}
	if info.WorkingDirectory != "" {
		sb.WriteString(fmt.Sprintf("cwd = %s\n", info.WorkingDirectory))
	}
	if info.Elevate {
		sb.WriteString("elevate = yes\n")
	}
	for key, value := range info.EnvVars {
		sb.WriteString(fmt.Sprintf("%s = %s\n", key, value))
	}

	sidecarName := fmt.Sprintf("%s.shim", info.Name)
	sidecarPath := filepath.Join(shimDir, sidecarName)

	err := os.WriteFile(sidecarPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", sidecarName, err)
	}

	return nil
}

func parseShimSidecar(sidecarPath string) (Info, error) {
	sidecarBytes, err := os.ReadFile(sidecarPath)
	if err != nil {
		return Info{}, fmt.Errorf("failed to read %s: %w", sidecarPath, err)
	}

	name := strings.TrimSuffix(filepath.Base(sidecarPath), ".shim")
	info := Info{Name: name, EnvVars: map[string]string{}}

	for line := range strings.SplitSeq(string(sidecarBytes), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "path":
			info.Target = value
		case "args":
			info.Args = value
		case "cwd":
			info.WorkingDirectory = value
		case "workdir":
			info.WorkingDirectory = value
		case "elevate":
			info.Elevate = value == "yes" || value == "true" || value == "1"
		case "runas":
			info.Elevate = value == "yes" || value == "true" || value == "1"
		case "":
			continue
		default:
			info.EnvVars[key] = value
		}
	}

	return info, nil
}
