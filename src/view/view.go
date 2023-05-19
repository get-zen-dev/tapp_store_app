package view

import (
	"constants"
	"environment"
	"shortQuestion"
	"style"
	"table"
	"theme"
	"time"

	k8 "k8sinterface"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ColumnTitleTitle          = "Title"
	ColumnTitleStatus         = "Status"
	ColumnTitleCurrentVersion = "Current"
	ColumnTitleLastVersion    = "Last"

	ColumnFlexTitle          = 3
	ColumnFlexStatus         = 3
	ColumnFlexCurrentVersion = 1
	ColumnFlexLastVersion    = 2

	ColumnMinSizeTitle          = 8
	ColumnMinSizeStatus         = 8
	ColumnMinSizeCurrentVersion = 10
	ColumnMinSizeLastVersion    = 10

	MinWidth  = 50
	MinHeight = 10
)

var (
	headers = []table.Column{
		{Title: ColumnTitleTitle, Width: ColumnMinSizeTitle, MinWidth: ColumnMinSizeTitle, Flex: ColumnFlexTitle},
		{Title: ColumnTitleStatus, Width: ColumnMinSizeStatus, MinWidth: ColumnMinSizeStatus, Flex: ColumnFlexStatus},
		{Title: ColumnTitleCurrentVersion, Width: ColumnMinSizeCurrentVersion, MinWidth: ColumnMinSizeCurrentVersion, Flex: ColumnFlexLastVersion},
		{Title: ColumnTitleLastVersion, Width: ColumnMinSizeLastVersion, MinWidth: ColumnMinSizeLastVersion, Flex: ColumnFlexLastVersion},
	}

	ratio = []int{
		ColumnFlexTitle,
		ColumnFlexStatus,
		ColumnFlexCurrentVersion,
		ColumnFlexLastVersion,
	}

	clientMicrok8s = k8.Microk8sClient{}
)

type Item struct {
	Title          string
	Status         string
	CurrentVersion string
	LastVersion    string
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
		item.CurrentVersion,
		item.LastVersion,
	}
}

func (i *Items) GetItems() []table.Row {
	return i.items
}

type Model struct {
	table    table.Model
	style    *style.Styles
	question shortQuestion.Question
}

func NewModel(data []table.Row) *Model {
	s := style.InitStyles(*theme.DefaultTheme)
	emptyState := "not found"
	m := Model{
		table: table.NewModel(s,
			constants.Dimensions{Width: MinWidth, Height: MinHeight},
			time.Now(),
			headers,
			data,
			"addons",
			&emptyState),
		style:    &s,
		question: shortQuestion.NewQuestionConcrete(),
	}
	_, err := environment.ReadFromConfig("domen")
	if err == nil {
		m.question.SetAnswered(true)
	}
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.question.Answered() {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.table.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height})
			m.table.SyncViewPortContent()
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				m.table.PrevItem()
			case "down", "j":
				m.table.NextItem()
			}
		}
		return m, nil
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.question.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height})
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				domen := m.question.Input().Value()
				environment.WriteInConfig("domen", domen)
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
		return m.table.View()
	}
	return m.question.View()
}
