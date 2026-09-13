package output

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func LogSuccess(format string, args ...any) {
	var output string

	if len(args) == 0 {
		output = activeTheme.SuccessColor.Sprint(format)
	} else if len(args) > 0 {
		output = activeTheme.SuccessColor.Sprintf(format, args...)
	}

	_, _ = fmt.Fprintln(os.Stderr, output)
}

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
			sb.WriteString(fmt.Sprintf("└── %s", lines[i]))
		} else {
			sb.WriteString(fmt.Sprintf("├── %s\n", lines[i]))
		}
	}

	_, _ = fmt.Fprintln(os.Stderr, activeTheme.ErrorColor.Sprint(sb.String()))
}

func LogWarning(format string, args ...any) {
	var output string

	if len(args) == 0 {
		output = activeTheme.WarningColor.Sprint(format)
	} else if len(args) > 0 {
		output = activeTheme.WarningColor.Sprintf(format, args...)
	}

	_, _ = fmt.Fprintln(os.Stderr, output)
}

func LogInfo(format string, args ...any) {
	var output string

	if len(args) == 0 {
		output = activeTheme.InfoColor.Sprint(format)
	} else if len(args) > 0 {
		output = activeTheme.InfoColor.Sprintf(format, args...)
	}

	_, _ = fmt.Fprintln(os.Stderr, output)
}

func LogOutput(format string, args ...any) {
	_, _ = fmt.Fprintln(os.Stdout, fmt.Sprintf(format, args...))
}
