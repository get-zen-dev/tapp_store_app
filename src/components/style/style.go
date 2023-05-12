package style

import (
	"common"
	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"theme"
)

type Styles struct {
	Common common.CommonStyles

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

	Section struct {
		ContainerPadding int
		ContainerStyle   lipgloss.Style
		SpinnerStyle     lipgloss.Style
		EmptyStateStyle  lipgloss.Style
		KeyStyle         lipgloss.Style
	}
}

func InitStyles(theme theme.Theme) Styles {
	var s Styles

	s.Common = common.BuildStyles(theme)

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

	s.Section.ContainerPadding = 1
	s.Section.ContainerStyle = lipgloss.NewStyle().Padding(0, s.Section.ContainerPadding)
	s.Section.SpinnerStyle = lipgloss.NewStyle().Padding(0, 1)
	s.Section.EmptyStateStyle = lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(1).
		MarginBottom(1)
	s.Section.KeyStyle = lipgloss.NewStyle().
		Foreground(theme.PrimaryText).
		Background(theme.SelectedBackground).
		Padding(0, 1)

	s.ListViewPort.PagerStyle = lipgloss.NewStyle().
		Height(common.ListPagerHeight).
		MaxHeight(common.ListPagerHeight).
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
	s.Table.SingleRuneTitleCellStyle = s.Table.TitleCellStyle.Copy().Width(common.SingleRuneWidth)
	s.Table.HeaderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.FaintBorder).
		BorderBottom(true)
	s.Table.RowStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.FaintBorder).
		BorderBottom(true)

	return s
}
