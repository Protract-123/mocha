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
		return err
	}

	separators := make([]string, len(headers))
	for i, h := range headers {
		separators[i] = strings.Repeat("-", len(h))
	}
	if _, err := fmt.Fprintln(tableWriter, strings.Join(separators, "\t")); err != nil {
		return err
	}

	for _, row := range rows {
		line := strings.Join(row, "\t")
		if _, err := fmt.Fprintln(tableWriter, line); err != nil {
			return err
		}
	}

	return tableWriter.Flush()
}
