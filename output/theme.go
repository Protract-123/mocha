package output

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Protract-123/mocha/config"
	"github.com/fatih/color"
)

type Theme struct {
	AccentColor  color.Color
	ErrorColor   color.Color
	WarningColor color.Color
	InfoColor    color.Color
}

var activeTheme = buildTheme(config.Default().Colors)
var initThemeOnce sync.Once

func InitTheme() {
	initThemeOnce.Do(func() {
		userTheme, err := resolveTheme(config.Current().Colors)
		if err != nil {
			LogWarning(fmt.Sprintf("warning: %v, using default theme", err))
			return
		}

		activeTheme = userTheme
	})
}

func buildTheme(cc config.ColorConfig) Theme {
	accent, _ := resolveColor(cc.AccentColor)
	errorColor, _ := resolveColor(cc.ErrorColor)
	warning, _ := resolveColor(cc.WarningColor)
	info, _ := resolveColor(cc.InfoColor)

	return Theme{
		AccentColor:  accent,
		ErrorColor:   errorColor,
		WarningColor: warning,
		InfoColor:    info,
	}
}

func resolveTheme(cc config.ColorConfig) (Theme, error) {
	accent, err := resolveColor(cc.AccentColor)
	if err != nil {
		return Theme{}, fmt.Errorf("theme.accent: %w", err)
	}

	errorColor, err := resolveColor(cc.ErrorColor)
	if err != nil {
		return Theme{}, fmt.Errorf("theme.error: %w", err)
	}

	warning, err := resolveColor(cc.WarningColor)
	if err != nil {
		return Theme{}, fmt.Errorf("theme.warning: %w", err)
	}

	info, err := resolveColor(cc.InfoColor)
	if err != nil {
		return Theme{}, fmt.Errorf("theme.info: %w", err)
	}

	return Theme{
		AccentColor:  accent,
		ErrorColor:   errorColor,
		WarningColor: warning,
		InfoColor:    info,
	}, nil
}

var colorNames = map[string]color.Attribute{
	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,

	"bright-red":     color.FgHiRed,
	"bright-green":   color.FgHiGreen,
	"bright-yellow":  color.FgHiYellow,
	"bright-blue":    color.FgHiBlue,
	"bright-magenta": color.FgHiMagenta,
	"bright-cyan":    color.FgHiCyan,
	"bright-white":   color.FgHiWhite,
}

func resolveColor(name string) (color.Color, error) {
	attr, ok := colorNames[strings.ToLower(name)]
	if !ok {
		return color.Color{}, fmt.Errorf("unknown color %q", name)
	}
	return *color.New(attr), nil
}
