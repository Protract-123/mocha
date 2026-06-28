package fileops

import (
	"fmt"
	"os/exec"
	"strings"
)

// TODO: consider native implementation / https://github.com/nyaosorg/go-windows-junction

func CreateJunction(targetDir string, junctionPath string) error {
	cmd := exec.Command("cmd", "/c", "mklink", "/j", junctionPath, targetDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create junction (output: %s): %w", strings.TrimSpace(string(out)), err)
	}
	return nil
}
