package view

import (
	"constants"
	"fmt"
	"listviewport"
	"strings"
	"theme"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	text = lipgloss.NewStyle().
		Foreground(theme.DefaultTheme.PrimaryText)
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Inherit(text)
	}()
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b).Inherit(text)
	}()
)

type OutputLog struct {
	content  string
	prev     *Model
	viewport viewport.Model
}

func NewOutputLog(content string, model *Model, width, height int) OutputLog {
	viewport := viewport.New(width, height)
	content = text.Render(content)
	viewport.SetContent(content)
	l := OutputLog{
		content:  content,
		prev:     model,
		viewport: viewport,
	}
	headerHeight := lipgloss.Height(l.headerView())
	footerHeight := lipgloss.Height(l.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	l.viewport.Width = width
	l.viewport.Height = height - verticalMarginHeight
	l.viewport.GotoBottom()
	return l
}

func (l OutputLog) Init() tea.Cmd {
	return nil
}

func (l OutputLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key := msg.String(); key == "ctrl+c" || key == "l" {
			return l.prev, nil
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(l.headerView())
		footerHeight := lipgloss.Height(l.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		l.viewport.Width = msg.Width
		l.viewport.Height = msg.Height - verticalMarginHeight
		l.prev.Width = msg.Width
		l.prev.Height = msg.Height
		l.prev.table.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height - constants.Keys.HeightShort - HeightMessage})
	}
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

func (l OutputLog) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, l.headerView(), l.viewport.View(), l.footerView())
}

func (l OutputLog) headerView() string {
	title := titleStyle.Render("Logger")
	line := strings.Repeat("─", listviewport.Max(0, l.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (l OutputLog) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", l.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", listviewport.Max(0, l.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
