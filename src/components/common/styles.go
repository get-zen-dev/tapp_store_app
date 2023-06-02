package common

import (
	"github.com/charmbracelet/lipgloss"
	"theme"
)

var (
	HeaderHeight       = 2
	FooterHeight       = 2
	TableHeaderHeight  = 2
	ListPagerHeight    = 2
)

type CommonStyles struct {
	MainTextStyle lipgloss.Style
	FooterStyle   lipgloss.Style
}

func BuildStyles(theme theme.Theme) CommonStyles {
	var s CommonStyles

	s.MainTextStyle = lipgloss.NewStyle().
		Foreground(theme.PrimaryText).
		Bold(true)
	s.FooterStyle = lipgloss.NewStyle().
		Height(FooterHeight - 1).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.PrimaryBorder)

	return s
}
