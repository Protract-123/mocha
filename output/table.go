package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintTable(headers []string, rows [][]string) error {
	tableWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	headerLine := strings.Join(headers, "\t")
	if _, err := fmt.Fprintln(tableWriter, headerLine); err != nil {
		return fmt.Errorf("failed to write header line: %w", err)
	}

	separators := make([]string, len(headers))
	for i, header := range headers {
		separators[i] = strings.Repeat("-", len(header))
	}
	if _, err := fmt.Fprintln(tableWriter, strings.Join(separators, "\t")); err != nil {
		return fmt.Errorf("failed to write separators: %w", err)
	}

	for _, row := range rows {
		line := strings.Join(row, "\t")
		if _, err := fmt.Fprintln(tableWriter, line); err != nil {
			return fmt.Errorf("failed to write a row: %w", err)
		}
	}

	return tableWriter.Flush()
}
