package view

import (
	"constants"
	env "environment"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	k8 "k8sinterface"
	"requests"
	"shortQuestion"
	"style"
	"table"
	"theme"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ColumnTitleTitle       = "Title"
	ColumnTitleStatus      = "Status"
	ColumnTitleDescription = "Description"

	ColumnFlexTitle       = 3
	ColumnFlexStatus      = 3
	ColumnFlexDescription = 3

	ColumnMinSizeTitle       = 8
	ColumnMinSizeStatus      = 8
	ColumnMinSizeDescription = 10

	MinWidth  = 50
	MinHeight = 10

	Installed = "✓"
	Deleted   = "✗"
)

var (
	headers = []table.Column{
		{Title: ColumnTitleTitle, Width: ColumnMinSizeTitle, MinWidth: ColumnMinSizeTitle, Flex: ColumnFlexTitle},
		{Title: ColumnTitleStatus, Width: ColumnMinSizeStatus, MinWidth: ColumnMinSizeStatus, Flex: ColumnFlexStatus},
		{Title: ColumnTitleDescription, Width: ColumnMinSizeDescription, MinWidth: ColumnMinSizeDescription, Flex: ColumnFlexDescription},
	}

	clientMicrok8s, _ = k8.GetInterfaceProvider("")
)

type Item struct {
	Title       string
	Status      string
	Description string
}

type Items struct {
	items []table.Row
}

func NewItems() *Items {
	return &Items{}
}

func (i *Items) Append(item *Item) {
	i.items = append(i.items, makeRow(item))
}

func makeRow(item *Item) []string {
	return []string{
		item.Title,
		item.Status,
		item.Description,
	}
}

func (i *Items) GetItems() []table.Row {
	return i.items
}

type Model struct {
	table    table.Model
	style    *style.Styles
	question shortQuestion.Question
	help     help.Model
}

func NewModel() (*Model, error) {
	env.SetUpEnv()
	list, err := requests.GetListAddons()
	if err != nil {
		return nil, err
	}
	items := NewItems()
	for _, v := range list.Value() {
		if v.Name == "common" {
			continue
		}
		info, err := clientMicrok8s.GetModuleInfo(v.Name)
		if err != nil {
			return nil, err
		}
		status := ""
		if info.IsEnabled {
			status = Installed
		} else {
			status = Deleted
		}
		items.Append(&Item{
			Title:       v.Name,
			Status:      status,
			Description: v.Path})
	}
	s := style.InitStyles(*theme.DefaultTheme)
	emptyState := "not found"
	m := Model{
		table: table.NewModel(s,
			constants.Dimensions{Width: MinWidth, Height: MinHeight},
			headers,
			items.GetItems(),
			&emptyState),
		style:    &s,
		question: shortQuestion.NewQuestionConcrete(),
		help:     help.New(),
	}
	_, err = env.ReadFromConfig("domen")
	if err == nil {
		m.question.SetAnswered(true)
	}
	return &m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

type Install struct {
	index int
}

type Delete struct {
	index int
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.question.Answered() {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width
			m.table.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height - constants.Keys.HeightShort})
			m.table.SyncViewPortContent()
		case Install:
			m.table.Rows[msg.index][1] = Installed
			m.table.SyncViewPortContent()
		case Delete:
			m.table.Rows[msg.index][1] = Deleted
			m.table.SyncViewPortContent()
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, constants.Keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, constants.Keys.Up):
				m.table.PrevItem()
			case key.Matches(msg, constants.Keys.Down):
				m.table.NextItem()
			case key.Matches(msg, constants.Keys.Install):
				index := m.table.GetCurrItem()
				name := m.table.Rows[index][0]
				return m, func() tea.Msg {
					clientMicrok8s.InstallModule(name)
					return Install{index}
				}
			case key.Matches(msg, constants.Keys.Delete):
				index := m.table.GetCurrItem()
				name := m.table.Rows[index][0]
				return m, func() tea.Msg {
					clientMicrok8s.RemoveModule(name)
					return Delete{index}
				}
			}
		}
		return m, nil
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.question.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height})
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, constants.Keys.QuitWithoutQ):
				return m, tea.Quit
			case key.Matches(msg, constants.Keys.Enter):
				domen := m.question.Input().Value()
				env.WriteInConfig("domen", domen)
				m.question.SetAnswered(true)
				m.table.SetDimensions(m.question.GetDimensions())
				m.table.SyncViewPortContent()
				return m, m.question.Input().Blur
			}
		}
		return m, m.question.Update(msg)
	}
}

func (m *Model) View() string {
	if m.question.Answered() {
		return lipgloss.JoinVertical(lipgloss.Left,
			m.table.View(), m.style.Common.FooterStyle.Width(m.help.Width).Render(m.help.View(constants.Keys)))
	}
	return m.question.View()
}
