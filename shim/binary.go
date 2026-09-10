package shim

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protract-123/mocha/fileops"
)

func InstallBinary(release Release, arch string, mochaDir string) error {
	tempDirectory := filepath.Join(mochaDir, "temp")
	if err := os.MkdirAll(tempDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDirectory)

	zipName := fmt.Sprintf("shim-%s.zip", arch)
	zipPath := filepath.Join(tempDirectory, zipName)

	downloadURL, err := url.JoinPath("https://github.com/ScoopInstaller/Shim/releases/download", release.Language, "v"+release.Version, zipName)
	if err != nil {
		return fmt.Errorf("failed to create download url: %w", err)
	}

	if err := fileops.DownloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("failed to download %s: %w", zipName, err)
	}

	if err := fileops.ExtractZip(zipPath, tempDirectory); err != nil {
		return fmt.Errorf("failed to extract %s: %w", zipName, err)
	}

	binaryBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe: %w", err)
	}

	checksumBytes, err := os.ReadFile(filepath.Join(tempDirectory, "shim.exe.sha256"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe.sha256: %w", err)
	}

	sumBytes := sha256.Sum256(binaryBytes)
	sum := hex.EncodeToString(sumBytes[:])

	checksum, _, _ := strings.Cut(strings.TrimSpace(string(checksumBytes)), " ")

	if sum != checksum {
		return fmt.Errorf("shim.exe hash does not match shim.exe.sha256")
	}

	binaryPath := filepath.Join(mochaDir, "shim.exe")
	if err := os.WriteFile(binaryPath, binaryBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write shim.exe to %s: %w", binaryPath, err)
	}

	return nil
}

const (
	peSubsystemGUI     = 2
	peSubsystemConsole = 3
	peHeaderOffset     = 0x3C // PE header offset location in a .exe
	peSubsystemOffset  = 0x5C // subsystem offset relative to PE header
)

func copyShimBinary(info Info, overrideSubsystem bool, mochaDir string) error {
	shimExe, err := os.ReadFile(filepath.Join(mochaDir, "shim.exe"))
	if err != nil {
		return fmt.Errorf("failed to read shim.exe: %w", err)
	}

	if overrideSubsystem {
		subsystem, err := getPESubsystem(info.Target)
		if err != nil {
			return fmt.Errorf("failed to get PE subsystem: %w", err)
		}

		if subsystem != peSubsystemGUI && subsystem != peSubsystemConsole {
			return fmt.Errorf("invalid subsystem %d", subsystem)
		}

		if err := patchPESubsystem(shimExe, subsystem); err != nil {
			return fmt.Errorf("failed to patch PE subsystem: %w", err)
		}
	}

	targetPath := filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.exe", info.Name))

	if err := os.WriteFile(targetPath, shimExe, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetPath, err)
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
