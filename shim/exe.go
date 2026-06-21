package shim

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateExeShim(name string, path string, mochaDir string) error {
	err := InitShims(mochaDir)
	if err != nil {
		return err
	}

	if !filepath.IsAbs(path) {
		return fmt.Errorf("%s is not an absolute path", path)
	}

	shimFilePath := filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.shim", name))
	shimFileContents := fmt.Sprintf("path = %s", path)

	err = os.WriteFile(shimFilePath, []byte(shimFileContents), os.ModePerm)
	if err != nil {
		return err
	}

	subsystem, err := GetPESubsystem(path)
	if err != nil {
		return err
	}

	if subsystem != 2 && subsystem != 3 {
		return fmt.Errorf("invalid subsystem %d", subsystem)
	}

	shimExe, err := os.Open(filepath.Join(mochaDir, "shim.exe"))
	if err != nil {
		return err
	}
	defer func(shimExe *os.File) {
		err := shimExe.Close()
		if err != nil {
			println(err.Error())
		}
	}(shimExe)

	destExe, err := os.Create(filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.exe", name)))
	if err != nil {
		return err
	}
	defer func(destExe *os.File) {
		err := destExe.Close()
		if err != nil && !errors.Is(err, os.ErrClosed) {
			println(err.Error())
		}
	}(destExe)

	_, err = io.Copy(destExe, shimExe)
	if err != nil {
		return err
	}

	err = destExe.Sync()
	if err != nil {
		return err
	}

	err = destExe.Close()
	if err != nil {
		return err
	}

	return SetPESubsystem(filepath.Join(mochaDir, "shims", fmt.Sprintf("%s.exe", name)), subsystem)
}

func SetPESubsystem(path string, subsystem uint16) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Read DOS header to find PE header offset
	// The offset to the PE header is at 0x3C in the DOS header
	var peOffset int32
	_, err = f.Seek(0x3C, io.SeekStart)
	if err != nil {
		return err
	}

	err = binary.Read(f, binary.LittleEndian, &peOffset)
	if err != nil {
		return err
	}

	// Skip PE signature (4) + FileHeader (20) = 24 bytes
	// Subsystem is at offset 68 into the OptionalHeader
	subsystemOffset := int64(peOffset) + 0x5c

	_, err = f.Seek(subsystemOffset, io.SeekStart)
	if err != nil {
		return err
	}

	return binary.Write(f, binary.LittleEndian, subsystem)
}

func GetPESubsystem(path string) (uint16, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)

	var peOffset int32
	_, err = f.Seek(0x3C, io.SeekStart)
	if err != nil {
		return 0, err
	}
	err = binary.Read(f, binary.LittleEndian, &peOffset)
	if err != nil {
		return 0, err
	}

	_, err = f.Seek(int64(peOffset)+0x5C, io.SeekStart)
	if err != nil {
		return 0, err
	}

	var subsystem uint16
	err = binary.Read(f, binary.LittleEndian, &subsystem)
	if err != nil {
		return 0, err
	}

	return subsystem, nil
}
