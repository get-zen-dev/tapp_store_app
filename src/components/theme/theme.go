package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	SelectedBackground lipgloss.AdaptiveColor
	SpinnerColor       lipgloss.AdaptiveColor
	PrimaryBorder      lipgloss.AdaptiveColor
	FaintBorder        lipgloss.AdaptiveColor
	BorderColor        lipgloss.AdaptiveColor
	SecondaryBorder    lipgloss.AdaptiveColor
	FaintText          lipgloss.AdaptiveColor
	PrimaryText        lipgloss.AdaptiveColor
	SecondaryText      lipgloss.AdaptiveColor
	InvertedText       lipgloss.AdaptiveColor
	SuccessText        lipgloss.AdaptiveColor
	WarningText        lipgloss.AdaptiveColor
}

var DefaultTheme = &Theme{
	PrimaryBorder:      lipgloss.AdaptiveColor{Light: "013", Dark: "008"},
	SecondaryBorder:    lipgloss.AdaptiveColor{Light: "008", Dark: "007"},
	SelectedBackground: lipgloss.AdaptiveColor{Light: "006", Dark: "#4B0082"},
	SpinnerColor:       lipgloss.AdaptiveColor{Light: "006", Dark: "#4B0082"},
	BorderColor:        lipgloss.AdaptiveColor{Light: "006", Dark: "36"},
	FaintBorder:        lipgloss.AdaptiveColor{Light: "254", Dark: "000"},
	PrimaryText:        lipgloss.AdaptiveColor{Light: "000", Dark: "015"},
	SecondaryText:      lipgloss.AdaptiveColor{Light: "244", Dark: "251"},
	FaintText:          lipgloss.AdaptiveColor{Light: "007", Dark: "245"},
	InvertedText:       lipgloss.AdaptiveColor{Light: "015", Dark: "236"},
	SuccessText:        lipgloss.AdaptiveColor{Light: "002", Dark: "002"},
	WarningText:        lipgloss.AdaptiveColor{Light: "001", Dark: "001"},
}
