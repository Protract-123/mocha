package output

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func LogError(providedError error) {
	if providedError == nil {
		return
	}

	var lines []string
	for providedError != nil {
		currentMsg := providedError.Error()
		nextErr := errors.Unwrap(providedError)

		if nextErr != nil {
			nextMsg := nextErr.Error()
			currentMsg = strings.TrimSuffix(currentMsg, ": "+nextMsg)
		}

		lines = append(lines, currentMsg)
		providedError = nextErr
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Error: %s\n", lines[0]))

	for i := 1; i < len(lines); i++ {
		if i == len(lines)-1 {
			sb.WriteString(fmt.Sprintf("└── Root Cause: %s", lines[i]))
		} else {
			sb.WriteString(fmt.Sprintf("├── %s\n", lines[i]))
		}
	}

	if _, err := activeTheme.ErrorColor.Fprintln(os.Stderr, sb.String()); err != nil {
		return
	}
}

func LogWarning(format string, args ...any) {
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	if _, err := activeTheme.WarningColor.Fprintln(os.Stderr, message); err != nil {
		return
	}
}

func LogInfo(format string, args ...any) {
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	if _, err := activeTheme.InfoColor.Fprintln(os.Stderr, message); err != nil {
		return
	}
}

func LogOutput(format string, args ...any) {
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	if _, err := color.New(color.FgWhite).Fprintln(os.Stdout, message); err != nil {
		return
	}
}
