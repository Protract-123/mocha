package app

import (
	"fmt"
	"runtime"
)

func GetSystemArch() (string, error) {
	cpuArch := runtime.GOARCH

	if cpuArch == "386" {
		return "32bit", nil
	} else if cpuArch == "amd64" {
		return "64bit", nil
	} else if cpuArch == "arm64" {
		return "arm64", nil
	}

	return "", fmt.Errorf("cpu architecture %q is unsupported", cpuArch)
}
