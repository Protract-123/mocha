package output

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type TableConfig struct {
	Spacing     int
	Alignments  []Alignment
	BorderStyle BorderStyle
}

type Alignment int

const (
	LeftAlign Alignment = iota
	RightAlign
	CenterAlign
)

type BorderStyle string

const (
	LightBorder  BorderStyle = "─"
	HeavyBorder  BorderStyle = "━"
	DoubleBorder BorderStyle = "═"
)

func PrintTable(headers []string, rows [][]string, config TableConfig) error {
	var sb strings.Builder
	widths := columnWidths(headers, rows, config.Spacing)

	sb.WriteString("\n")

	for i, header := range headers {
		sb.WriteString(alignString(header, widths[i], config.Alignments[i]))
	}
	sb.WriteString("\n")

	for _, width := range widths {
		sb.WriteString(strings.Repeat(string(config.BorderStyle), width))
	}
	sb.WriteString("\n")

	for _, row := range rows {
		for j, item := range row {
			sb.WriteString(alignString(item, widths[j], config.Alignments[j]))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	_, err := fmt.Fprint(os.Stdout, sb.String())
	return err
}

func alignString(str string, columnWidth int, alignment Alignment) string {
	padding := columnWidth - utf8.RuneCountInString(str)

	switch alignment {
	case LeftAlign:
		return str + strings.Repeat(" ", padding)
	case RightAlign:
		return strings.Repeat(" ", padding) + str
	case CenterAlign:
		leftPadding := padding / 2
		rightPadding := padding - leftPadding
		return strings.Repeat(" ", leftPadding) + str + strings.Repeat(" ", rightPadding)
	default:
		return ""
	}
}

func columnWidths(headers []string, rows [][]string, spacing int) []int {
	headerCount := len(headers)
	widths := make([]int, headerCount)

	for i, header := range headers {
		columnWidth := 0

		for _, row := range rows {
			if width := utf8.RuneCountInString(row[i]); width > columnWidth {
				columnWidth = width
			}
		}

		if width := utf8.RuneCountInString(header); width > columnWidth {
			columnWidth = width
		}

		if i != headerCount-1 {
			columnWidth += spacing
		}

		widths[i] = columnWidth
	}

	return widths
}
