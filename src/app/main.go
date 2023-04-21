package main

import (
	"fmt"
	"os"

	env "environment"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"requests"
)

var columnStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	title string
	desc  string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func InitModel() (*model, error) {
	response, err := requests.GetListAddons()
	if err != nil {
		return nil, err
	}
	items := []list.Item{}
	for _, v := range response.Value() {
		items = append(items, Item{title: v.Name, desc: v.Path})
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "List addons"
	return &m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := columnStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return columnStyle.Render(m.list.View())
}

func main() {
	env.SetUpEnv()
	m, err := InitModel()
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(*m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
}
