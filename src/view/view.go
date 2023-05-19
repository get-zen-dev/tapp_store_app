package view

import (
	"constants"
	"style"
	"table"
	"theme"
	"time"

	k8 "k8sinterface"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	COLUMN_KEY_TITLE           = "TITLE"
	COLUMN_KEY_STATUS          = "STATUS"
	COLUMN_KEY_CURRENT_VERSION = "CURRENT"
	COLUMN_KEY_LAST_VERSION    = "LAST"

	COLUMN_INDEX_TITLE           = 0
	COLUMN_INDEX_STATUS          = 1
	COLUMN_INDEX_CURRENT_VERSION = 2
	COLUMN_INDEX_LAST_VERSION    = 3

	COLUMN_TITLE_TITLE           = "Title"
	COLUMN_TITLE_STATUS          = "Status"
	COLUMN_TITLE_CURRENT_VERSION = "Current"
	COLUMN_TITLE_LAST_VERSION    = "Last"

	COLUMN_FLEX_TITLE           = 3
	COLUMN_FLEX_STATUS          = 3
	COLUMN_FLEX_CURRENT_VERSION = 1
	COLUMN_FLEX_LAST_VERSION    = 2

	COLUMN_MIN_SIZE_TITLE           = 8
	COLUMN_MIN_SIZE_STATUS          = 8
	COLUMN_MIN_SIZE_CURRENT_VERSION = 10
	COLUMN_MIN_SIZE_LAST_VERSION    = 10

	MIN_WIDTH  = 50
	MIN_HEIGHT = 10

	MIN_HEIGHT_HELP = 6
)

var (
	headers = []table.Column{
		{Title: COLUMN_TITLE_TITLE, Width: COLUMN_MIN_SIZE_TITLE, MinWidth: COLUMN_MIN_SIZE_TITLE, Flex: COLUMN_FLEX_TITLE},
		{Title: COLUMN_TITLE_STATUS, Width: COLUMN_MIN_SIZE_STATUS, MinWidth: COLUMN_MIN_SIZE_STATUS, Flex: COLUMN_FLEX_STATUS},
		{Title: COLUMN_TITLE_CURRENT_VERSION, Width: COLUMN_MIN_SIZE_CURRENT_VERSION, MinWidth: COLUMN_MIN_SIZE_CURRENT_VERSION, Flex: COLUMN_FLEX_LAST_VERSION},
		{Title: COLUMN_TITLE_LAST_VERSION, Width: COLUMN_MIN_SIZE_LAST_VERSION, MinWidth: COLUMN_MIN_SIZE_LAST_VERSION, Flex: COLUMN_FLEX_LAST_VERSION},
	}

	ratio = []int{
		COLUMN_FLEX_TITLE,
		COLUMN_FLEX_STATUS,
		COLUMN_FLEX_CURRENT_VERSION,
		COLUMN_FLEX_LAST_VERSION,
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
	table table.Model
	style *style.Styles
}

func NewModel(data []table.Row) *Model {
	s := style.InitStyles(*theme.DefaultTheme)
	emptyState := "not found"
	m := Model{
		table: table.NewModel(s,
			constants.Dimensions{Width: MIN_WIDTH, Height: MIN_HEIGHT},
			time.Now(),
			headers,
			data,
			"addons",
			&emptyState),
		style: &s,
	}
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		return m, nil
	}

	return m, nil
}

func (m *Model) View() string {
	return m.table.View()
}
