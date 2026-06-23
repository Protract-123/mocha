package output

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
	_, err := fmt.Fprintf(os.Stderr, "%s%s%s\n", AnsiRed, sb.String(), AnsiReset)
	if err != nil {
		return
	}
}

func LogWarning(message string) {
	_, err := fmt.Fprintf(os.Stderr, "%sWarning: %s%s\n", AnsiYellow, message, AnsiReset)
	if err != nil {
		return
	}
}

func LogInfo(message string) {
	_, err := fmt.Fprintf(os.Stderr, "%s%s%s\n", AnsiBlue, message, AnsiReset)
	if err != nil {
		return
	}
}

func LogOutput(message string) {
	_, err := fmt.Printf("%s%s%s\n", AnsiWhite, message, AnsiReset)
	if err != nil {
		return
	}
}
