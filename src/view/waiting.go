package view

import (
	"style"
	"theme"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Waiting struct {
	spinner  spinner.Model
	style    style.StylesWaiting
	width    int
	height   int
	quitting bool
	sent     bool
	err      error
	command  func() bool
}

var points = spinner.Spinner{
	Frames: []string{"∙∙∙∙∙", "●∙∙∙∙", "∙●∙∙∙", "∙∙●∙∙", "∙∙∙●∙", "∙∙∙∙●"},
	FPS:    time.Second / 7,
}

func NewModelWaiting(fn func() bool) Waiting {
	spin := spinner.New()
	spin.Spinner = points
	s := style.InitStylesWaiting(*theme.DefaultTheme)
	spin.Style = s.Spinner
	return Waiting{spinner: spin, style: s, command: fn}
}

func (m Waiting) Init() tea.Cmd {
	return m.spinner.Tick
}

type Answer struct {
	value bool
}

func (m Waiting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case errMsg:
		m.err = msg
		return m, nil
	case Answer:
		kub, _ := NewModelTable()
		return kub.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	if !m.sent {
		fn := func() tea.Msg {
			return Answer{m.command()}
		}
		cmds = append(cmds, fn)
		m.sent = true
	}
	return m, tea.Batch(cmds...)
}

func (m Waiting) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.spinner.View(),
			m.style.Text.Render("Wait. Kubernetes launches"),
		),
	)
}
