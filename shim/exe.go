package shim

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

const (
	peSubsystemGUI     = 2
	peSubsystemConsole = 3
	peHeaderOffset     = 0x3C // PE header offset location in a .exe
	peSubsystemOffset  = 0x5C // subsystem offset relative to PE header
)

func CreateExeShim(name string, execPath string, mochaDir string) error {
	if err := InitShimBinary(mochaDir); err != nil {
		return fmt.Errorf("failed to init shim.exe: %w", err)
	}

	if !filepath.IsAbs(execPath) {
		return fmt.Errorf("%s is not an absolute path", execPath)
	}

	if err := createShimFile(name, execPath, mochaDir); err != nil {
		return fmt.Errorf("failed to create .shim file: %w", err)
	}

	if err := createShimExe(name, execPath, mochaDir); err != nil {
		return fmt.Errorf("failed to create shim .exe file: %w", err)
	}

	return nil
}

func createShimExe(name string, execPath string, mochaDir string) error {
	subsystem, err := getPESubsystem(execPath)
	if err != nil {
		return fmt.Errorf("failed to get PE subsystem: %w", err)
	}

	if subsystem != peSubsystemGUI && subsystem != peSubsystemConsole {
		return fmt.Errorf("invalid subsystem %d", subsystem)
	}

	shimExe, err := os.ReadFile(filepath.Join(mochaDir, "shim.exe"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe: %w", err)
	}

	if err := patchPESubsystem(shimExe, subsystem); err != nil {
		return fmt.Errorf("failed to patch PE subsystem: %w", err)
	}

	targetPath := filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.exe", name))

	if err := os.WriteFile(targetPath, shimExe, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetPath, err)
	}

	return nil
}

func createShimFile(name string, execPath string, mochaDir string) error {
	shimFileName := fmt.Sprintf("%s.shim", name)
	shimFilePath := filepath.Join(mochaDir, "shims", shimFileName)
	shimFileContents := fmt.Sprintf("path = %s", execPath)

	if err := os.WriteFile(shimFilePath, []byte(shimFileContents), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write %s: %w", shimFileName, err)
	}

	return nil
}

func patchPESubsystem(exeBytes []byte, subsystem uint16) error {
	subsystemOffset, err := getSubsystemOffset(exeBytes)
	if err != nil {
		return fmt.Errorf("failed to get subsystem offset: %w", err)
	}

	if len(exeBytes) < int(subsystemOffset)+2 {
		return fmt.Errorf("invalid exe: PE header out of bounds")
	}

	binary.LittleEndian.PutUint16(exeBytes[subsystemOffset:subsystemOffset+2], subsystem)

	return nil
}

func getPESubsystem(path string) (uint16, error) {
	exeFile, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", path, err)
	}

	subsystemOffset, err := getSubsystemOffset(exeFile)
	if err != nil {
		return 0, fmt.Errorf("failed to get subsystem offset: %w", err)
	}

	if len(exeFile) < int(subsystemOffset)+2 {
		return 0, fmt.Errorf("invalid exe: PE header out of bounds")
	}

	subsystem := binary.LittleEndian.Uint16(exeFile[subsystemOffset : subsystemOffset+2])
	return subsystem, nil
}

func getSubsystemOffset(exeBytes []byte) (uint32, error) {
	if len(exeBytes) < peHeaderOffset+4 {
		return 0, fmt.Errorf("invalid exe: file too small")
	}

	peOffsetBytes := exeBytes[peHeaderOffset : peHeaderOffset+4]
	subsystemOffset := binary.LittleEndian.Uint32(peOffsetBytes) + uint32(peSubsystemOffset)

	return subsystemOffset, nil
}
