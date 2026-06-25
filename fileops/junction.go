package fileops

import (
	"fmt"
	"os/exec"
	"strings"
)

// TODO: consider native implementation / https://github.com/nyaosorg/go-windows-junction

func CreateJunction(targetDir string, junctionPath string) error {
	cmd := exec.Command("cmd", "/c", "mklink", "/j", junctionPath, targetDir)

	out, err := cmd.CombinedOutput()
	if err != nil {
		println(cmd.String())
		return fmt.Errorf("failed to create junction. Output: %s: %w", strings.TrimSpace(string(out)), err)
	}

	return nil
}
