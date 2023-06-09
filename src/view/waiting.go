package view

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	he "handleException"
	k8 "k8sinterface"
	"style"
	"theme"
	"time"
)

var (
	KubernetesLaunch = "Wait. Kubernetes launches"
)

type errMsg error

type Waiting struct {
	spinner        spinner.Model
	style          style.StylesWaiting
	text           string
	width          int
	height         int
	quitting       bool
	sent           bool
	err            error
	clientMicrok8s k8.KuberInterface
}

var points = spinner.Spinner{
	Frames: []string{"∙∙∙∙∙", "●∙∙∙∙", "∙●∙∙∙", "∙∙●∙∙", "∙∙∙●∙", "∙∙∙∙●"},
	FPS:    time.Second / 7,
}

func NewModelWaiting(clientMicrok8s k8.KuberInterface, text string) Waiting {
	spin := spinner.New()
	spin.Spinner = points
	s := style.InitStylesWaiting(*theme.DefaultTheme)
	spin.Style = s.Spinner
	return Waiting{spinner: spin, style: s, text: text, clientMicrok8s: clientMicrok8s}
}

func (m Waiting) Init() tea.Cmd {
	return m.spinner.Tick
}

type Next struct {
	value error
}

func (m Waiting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case errMsg:
		m.err = msg
		return m, nil
	case Next:
		he.PrintErrorIfNotNil(msg.value)
		next, _ := NewModelTable(m.clientMicrok8s, m.width, m.height)
		return next, nil
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
			command := func() error {
				return m.clientMicrok8s.Start()
			}
			return Next{command()}
		}
		cmds = append(cmds, fn)
		m.sent = true
	}
	return m, tea.Batch(cmds...)
}

func (m Waiting) View() string {
	he.PrintErrorIfNotNil(m.err)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.spinner.View(),
			m.style.Text.Render(m.text),
		),
	)
}
