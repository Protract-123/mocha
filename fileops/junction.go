// The CreateJunction function below is adapted from
// https://github.com/nyaosorg/go-windows-junction, used under the MIT License below.
//
// MIT License
//
// Copyright (c) 2020 zetamatta
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// Modifications: wrapped errors with %w and made error messages descriptive.

package fileops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

// Value for FSCTL_SET_REPARSE_POINT can be found in the following locations:
// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/5c03c9d6-15de-48a2-9835-8fb37f8a79d8
// https://go.dev/src/internal/syscall/windows/reparse_windows.go

const kFSCTL_SET_REPARSE_POINT = 589988

func CreateJunction(target, mountPt string) error {
	_target, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %q: %w", target, err)
	}
	_mountPt, err := windows.UTF16PtrFromString(mountPt)
	if err != nil {
		return fmt.Errorf("failed to make UTF16Ptr from %q: %w", mountPt, err)
	}

	err = os.Mkdir(mountPt, 0777)
	if err != nil {
		return fmt.Errorf("failed to make target dir %q: %w", mountPt, err)
	}
	ok := false
	defer func() {
		if !ok {
			os.Remove(mountPt)
		}
	}()

	handle, err := windows.CreateFile(_mountPt,
		windows.GENERIC_WRITE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0)
	if err != nil {
		return fmt.Errorf("failed to create windows file handle for %q: %w", mountPt, err)
	}
	defer windows.CloseHandle(handle)

	rp := winio.ReparsePoint{
		Target:       _target,
		IsMountPoint: true,
	}

	data := winio.EncodeReparsePoint(&rp)

	var size uint32

	err = windows.DeviceIoControl(
		handle,
		kFSCTL_SET_REPARSE_POINT,
		&data[0],
		uint32(len(data)),
		nil,
		0,
		&size,
		nil)

	if err != nil {
		return fmt.Errorf("failed to set reparse point: %w", err)
	}
	ok = true
	return nil
}
