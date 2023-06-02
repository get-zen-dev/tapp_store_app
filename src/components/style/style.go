package style

import (
	bbHelp "github.com/charmbracelet/bubbles/help"
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
	FooterStyle   lipgloss.Style
}

func BuildStyles(theme theme.Theme) CommonStyles {
	var s CommonStyles

	s.FooterStyle = lipgloss.NewStyle().
		Height(FooterHeight - 1).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.PrimaryBorder)

	return s
}

type Styles struct {
	Common CommonStyles

	Help struct {
		Text         lipgloss.Style
		KeyText      lipgloss.Style
		BubbleStyles bbHelp.Styles
	}
	ListViewPort struct {
		PagerHeight int
		PagerStyle  lipgloss.Style
	}
	Table struct {
		CellStyle                lipgloss.Style
		SelectedCellStyle        lipgloss.Style
		TitleCellStyle           lipgloss.Style
		SingleRuneTitleCellStyle lipgloss.Style
		HeaderStyle              lipgloss.Style
		RowStyle                 lipgloss.Style
	}
}

func InitStyles(theme theme.Theme) Styles {
	var s Styles

	s.Common = BuildStyles(theme)

	s.Help.Text = lipgloss.NewStyle().Foreground(theme.SecondaryText)
	s.Help.KeyText = lipgloss.NewStyle().Foreground(theme.PrimaryText)
	s.Help.BubbleStyles = bbHelp.Styles{
		ShortDesc:      s.Help.Text.Copy().Foreground(theme.FaintText),
		FullDesc:       s.Help.Text.Copy().Foreground(theme.FaintText),
		ShortSeparator: s.Help.Text.Copy().Foreground(theme.SecondaryBorder),
		FullSeparator:  s.Help.Text.Copy(),
		FullKey:        s.Help.KeyText.Copy(),
		ShortKey:       s.Help.KeyText.Copy(),
		Ellipsis:       s.Help.Text.Copy(),
	}

	s.ListViewPort.PagerStyle = lipgloss.NewStyle().
		Height(ListPagerHeight).
		MaxHeight(ListPagerHeight).
		PaddingTop(1).
		Foreground(theme.FaintText)

	s.Table.CellStyle = lipgloss.NewStyle().PaddingLeft(1).
		PaddingRight(1).
		MaxHeight(1)
	s.Table.SelectedCellStyle = s.Table.CellStyle.Copy().
		Background(theme.SelectedBackground)
	s.Table.TitleCellStyle = s.Table.CellStyle.Copy().
		Bold(true).
		Foreground(theme.PrimaryText)
	s.Table.HeaderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.FaintBorder).
		BorderBottom(true)
	s.Table.RowStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(theme.FaintBorder).
		BorderBottom(true)

	return s
}

type StylesWaiting struct {
	Spinner lipgloss.Style
	Text    lipgloss.Style
}

func InitStylesWaiting(theme theme.Theme) StylesWaiting {
	var style StylesWaiting
	style.Spinner = lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		Foreground(theme.SpinnerColor).
		BorderBottom(true)
	style.Text = lipgloss.NewStyle().
		Foreground(theme.PrimaryText)
	return style
}

type StylesQuestion struct {
	BorderColor lipgloss.AdaptiveColor
	InputField  lipgloss.Style
}

func InitStylesQuestion(theme theme.Theme) StylesQuestion {
	s := StylesQuestion{}
	s.BorderColor = theme.BorderColor
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1)
	return s
}
