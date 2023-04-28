package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	COLUMN_KEY_ID              = "ID"
	COLUMN_KEY_IMAGE           = "IMAGE"
	COLUMN_KEY_TITLE           = "TITLE"
	COLUMN_KEY_STATUS          = "STATUS"
	COLUMN_KEY_CURRENT_VERSION = "CURRENT_VERSION"
	COLUMN_KEY_LAST_VERSION    = "LAST_VERSION"

	COLUMN_TITLE_ID              = "id"
	COLUMN_TITLE_IMAGE           = "image"
	COLUMN_TITLE_TITLE           = "title"
	COLUMN_TITLE_STATUS          = "status"
	COLUMN_TITLE_CURRENT_VERSION = "current version"
	COLUMN_TITLE_LAST_VERSION    = "last version"

	COLUMN_FLEX_ID              = 1
	COLUMN_FLEX_IMAGE           = 3
	COLUMN_FLEX_TITLE           = 6
	COLUMN_FLEX_STATUS          = 6
	COLUMN_FLEX_CURRENT_VERSION = 3
	COLUMN_FLEX_LAST_VERSION    = 3

	MARGI_RIGHT = 0
	MARGI_LEFT  = 0

	colorWater = "#22BEDA"
	colorNeon  = "#20CF97"

	PER_PAGE = 7
)

var (
	styleBase = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorWater)).
		BorderForeground(lipgloss.Color(colorNeon)).
		Align(lipgloss.Right)
)

type Item struct {
	Image          string
	Title          string
	Status         string
	CurrentVersion string
	LastVersion    string
}

type Items struct {
	slice []table.Row
}

func NewItems(newSlice []table.Row) *Items {
	return &Items{newSlice}
}

func (i *Items) Append(item Item) {
	i.slice = append(i.slice, makeRow(fmt.Sprintf("%v", len(i.slice)+1), item.Image, item.Title, item.Status, item.CurrentVersion, item.LastVersion))
}

func (i *Items) GetRows() []table.Row {
	return i.slice
}

func makeRow(ID, image, title, status, currentVersion, lastVersion string) table.Row {
	return table.NewRow(table.RowData{
		COLUMN_KEY_ID:              ID,
		COLUMN_KEY_IMAGE:           image,
		COLUMN_KEY_TITLE:           title,
		COLUMN_KEY_STATUS:          status,
		COLUMN_KEY_CURRENT_VERSION: currentVersion,
		COLUMN_KEY_LAST_VERSION:    lastVersion,
	})
}

type Model struct {
	addonTable      table.Model
	filterTextInput textinput.Model
	totalWidth      int
	totalMargin     int
}

func NewModel(data []table.Row) Model {
	return Model{
		addonTable: table.New([]table.Column{
			table.NewFlexColumn(COLUMN_KEY_ID, COLUMN_TITLE_ID, COLUMN_FLEX_ID),
			table.NewFlexColumn(COLUMN_KEY_IMAGE, COLUMN_TITLE_IMAGE, COLUMN_FLEX_IMAGE),
			table.NewFlexColumn(COLUMN_KEY_TITLE, COLUMN_TITLE_TITLE, COLUMN_FLEX_TITLE),
			table.NewFlexColumn(COLUMN_KEY_STATUS, COLUMN_TITLE_STATUS, COLUMN_FLEX_STATUS),
			table.NewFlexColumn(COLUMN_KEY_CURRENT_VERSION, COLUMN_TITLE_CURRENT_VERSION, COLUMN_FLEX_CURRENT_VERSION),
			table.NewFlexColumn(COLUMN_KEY_LAST_VERSION, COLUMN_TITLE_LAST_VERSION, COLUMN_FLEX_LAST_VERSION),
		}).WithRows(data).
			Filtered(true).
			BorderRounded().
			WithPageSize(PER_PAGE).
			SortByAsc(COLUMN_KEY_TITLE).
			WithBaseStyle(styleBase).
			Focused(true),
		totalMargin:     MARGI_LEFT + MARGI_RIGHT,
		filterTextInput: textinput.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.addonTable, cmd = m.addonTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) recalculateTable() {
	m.addonTable = m.addonTable.WithTargetWidth(m.totalWidth - m.totalMargin)
}

func (m Model) View() string {
	view := lipgloss.JoinVertical(lipgloss.Right, m.addonTable.View()) + "\n"
	return lipgloss.NewStyle().Render(view)
}
