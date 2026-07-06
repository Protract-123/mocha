// The CreateShortcut function below is adapted from
// https://github.com/jxeng/shortcut, used under the MIT License below.
// MIT License
//
// Copyright (c) 2022 jxeng
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
// Modifications:
//    rename Create(shortcut Shortcut) to CreateShortcut(shortcut Shortcut)
//    handle CoInitializeEx error and defer CoUninitialize
//    add defer idispatch.Release()
//    wrapped errors with %w and made error messages descriptive

package fileops

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type Shortcut struct {
	// Shortcut (.lnk file) path
	ShortcutPath string
	// Shortcut target: a file path or a website
	Target string
	// Shortcut icon path, default: "%SystemRoot%\\System32\\SHELL32.dll,0"
	IconLocation string
	// Arguments of shortcut
	Arguments string
	// Description of shortcut
	Description string
	// Hotkey of shortcut
	Hotkey string
	// WindowStyle, "1"(default) for default size and location; "3" for maximized window; "7" for minimized window
	WindowStyle string
	// Working directory of shortcut
	WorkingDirectory string
}

func CreateShortcut(shortcut Shortcut) error {
	if shortcut.IconLocation == "" {
		shortcut.IconLocation = "%SystemRoot%\\System32\\SHELL32.dll,0"
	}

	if shortcut.WindowStyle == "" {
		shortcut.WindowStyle = "1"
	}

	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY); err != nil {
		return fmt.Errorf("failed to initialize COM: %w", err)
	}
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("failed to create WScript.Shell: %w", err)
	}
	defer oleShellObject.Release()

	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to get wshell: %w", err)
	}
	defer wshell.Release()

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcut.ShortcutPath)
	if err != nil {
		return fmt.Errorf("failed to create shortcut: %w", err)
	}
	idispatch := cs.ToIDispatch()
	defer idispatch.Release()

	if _, err := oleutil.PutProperty(idispatch, "IconLocation", shortcut.IconLocation); err != nil {
		return fmt.Errorf("failed to set icon location: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "TargetPath", shortcut.Target); err != nil {
		return fmt.Errorf("failed to set target path: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "Arguments", shortcut.Arguments); err != nil {
		return fmt.Errorf("failed to set arguments: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "Description", shortcut.Description); err != nil {
		return fmt.Errorf("failed to set description: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "Hotkey", shortcut.Hotkey); err != nil {
		return fmt.Errorf("failed to set hotkey: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "WindowStyle", shortcut.WindowStyle); err != nil {
		return fmt.Errorf("failed to set window style: %w", err)
	}
	if _, err := oleutil.PutProperty(idispatch, "WorkingDirectory", shortcut.WorkingDirectory); err != nil {
		return fmt.Errorf("failed to set working directory: %w", err)
	}
	if _, err := oleutil.CallMethod(idispatch, "Save"); err != nil {
		return fmt.Errorf("failed to save shortcut: %w", err)
	}
	return nil
}
