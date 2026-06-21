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
	err := InitShims(mochaDir)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(execPath) {
		return fmt.Errorf("%s is not an absolute path", execPath)
	}

	err = createShimFile(name, execPath, mochaDir)
	if err != nil {
		return fmt.Errorf("failed to create .shim file: %w", err)
	}

	err = createShimExe(name, execPath, mochaDir)
	if err != nil {
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

	err = patchPESubsystem(shimExe, subsystem)
	if err != nil {
		return fmt.Errorf("failed to patch PE subsystem: %w", err)
	}

	targetPath := filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.exe", name))
	err = os.WriteFile(targetPath, shimExe, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", targetPath, err)
	}

	return nil
}

func createShimFile(name string, execPath string, mochaDir string) error {
	shimFileName := fmt.Sprintf("%s.shim", name)
	shimFilePath := filepath.Join(mochaDir, "shims", shimFileName)
	shimFileContents := fmt.Sprintf("path = %s", execPath)

	err := os.WriteFile(shimFilePath, []byte(shimFileContents), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", shimFileName, err)
	}

	return nil
}

func patchPESubsystem(exeFile []byte, subsystem uint16) error {
	if len(exeFile) < peHeaderOffset+4 {
		return fmt.Errorf("invalid exe: file too small")
	}

	peOffsetBytes := exeFile[peHeaderOffset : peHeaderOffset+4]
	peOffset := binary.LittleEndian.Uint32(peOffsetBytes)

	subsystemOffset := peOffset + uint32(peSubsystemOffset)
	if len(exeFile) < int(subsystemOffset)+2 {
		return fmt.Errorf("invalid exe: PE header out of bounds")
	}

	binary.LittleEndian.PutUint16(exeFile[subsystemOffset:subsystemOffset+2], subsystem)

	return nil
}

func getPESubsystem(path string) (uint16, error) {
	exeFile, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", path, err)
	}

	if len(exeFile) < peHeaderOffset+4 {
		return 0, fmt.Errorf("invalid exe: file too small")
	}

	peOffsetBytes := exeFile[peHeaderOffset : peHeaderOffset+4]
	peOffset := binary.LittleEndian.Uint32(peOffsetBytes)

	subsystemOffset := peOffset + uint32(peSubsystemOffset)

	if len(exeFile) < int(subsystemOffset)+2 {
		return 0, fmt.Errorf("invalid exe: PE header out of bounds")
	}

	subsystem := binary.LittleEndian.Uint16(exeFile[subsystemOffset : subsystemOffset+2])
	return subsystem, nil
}
