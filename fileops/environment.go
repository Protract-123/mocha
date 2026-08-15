package fileops

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func SetEnvironmentVariable(name string, value string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open environment registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(name, value); err != nil {
		return fmt.Errorf("failed to set environment registry value: %w", err)
	}

	return nil
}

func RemoveEnvironmentVariable(name string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open environment registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(name); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return fmt.Errorf("failed to delete environment registry value: %w", err)
	}

	return nil
}

func GetEnvironmentVariable(name string) (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open environment registry key: %w", err)
	}
	defer key.Close()

	val, _, err := key.GetStringValue(name)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get environment registry value: %w", err)
	}

	return val, nil
}

func PropagateEnvironment() error {
	// Value for WM_SETTINGCHANGE can be found in the following location:
	// https://learn.microsoft.com/en-us/windows/win32/winmsg/wm-settingchange
	const WM_SETTINGCHANGE = 0x001A

	// Value for HWND_BROADCAST can be found in the following location:
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendnotifymessagea
	const HWND_BROADCAST = 0xFFFF

	// Value for SMTO_ABORTIFHUNG can be found in the following location:
	// see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagetimeouta
	const SMTO_ABORTIFHUNG = 0x0002

	const timeoutMs = 5000

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	envStr, err := syscall.UTF16PtrFromString("Environment")
	if err != nil {
		return fmt.Errorf("failed to encode environment string: %w", err)
	}

	sendMessageTimeout.Call(
		HWND_BROADCAST,
		WM_SETTINGCHANGE,
		0,
		uintptr(unsafe.Pointer(envStr)),
		SMTO_ABORTIFHUNG,
		timeoutMs,
		0,
	)

	return nil
}
