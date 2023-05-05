package view

import (
	"fmt"
	"unicode"

	k8 "k8sinterface"

	"github.com/76creates/stickers/flexbox"
	"github.com/76creates/stickers/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	COLUMN_KEY_TITLE           = "TITLE"
	COLUMN_KEY_STATUS          = "STATUS"
	COLUMN_KEY_CURRENT_VERSION = "CURRENT_VERSION"
	COLUMN_KEY_LAST_VERSION    = "LAST_VERSION"

	COLUMN_INDEX_TITLE           = 0
	COLUMN_INDEX_STATUS          = 1
	COLUMN_INDEX_CURRENT_VERSION = 2
	COLUMN_INDEX_LAST_VERSION    = 3

	COLUMN_TITLE_TITLE           = "title"
	COLUMN_TITLE_STATUS          = "status"
	COLUMN_TITLE_CURRENT_VERSION = "current version"
	COLUMN_TITLE_LAST_VERSION    = "last version"

	COLUMN_FLEX_TITLE           = 6
	COLUMN_FLEX_STATUS          = 6
	COLUMN_FLEX_CURRENT_VERSION = 3
	COLUMN_FLEX_LAST_VERSION    = 3

	COLUMN_MIN_SIZE_TITLE           = 6
	COLUMN_MIN_SIZE_STATUS          = 6
	COLUMN_MIN_SIZE_CURRENT_VERSION = 3
	COLUMN_MIN_SIZE_LAST_VERSION    = 3

	MARGIN_RIGHT = 0
	MARGIN_LEFT  = 0

	MIN_HEIGHT_HELP = 6

	help = `
move: ←, ↑, →, ↓ | install: i, delete: d, update: u | quit: ctrl+c`
)

var (
	headers = []string{
		COLUMN_KEY_TITLE,
		COLUMN_KEY_STATUS,
		COLUMN_KEY_CURRENT_VERSION,
		COLUMN_KEY_LAST_VERSION,
	}

	ratio = []int{
		COLUMN_FLEX_TITLE,
		COLUMN_FLEX_STATUS,
		COLUMN_FLEX_CURRENT_VERSION,
		COLUMN_FLEX_LAST_VERSION,
	}

	minSize = []int{
		COLUMN_MIN_SIZE_TITLE,
		COLUMN_MIN_SIZE_STATUS,
		COLUMN_MIN_SIZE_CURRENT_VERSION,
		COLUMN_MIN_SIZE_LAST_VERSION,
	}

	clientMicrok8s = k8.Microk8sClient{}

	selectedValue = "\nselect something with spacebar or enter"
)

type Item struct {
	Title          string
	Status         string
	CurrentVersion string
	LastVersion    string
}

type Items struct {
	items [][]string
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

func (i *Items) GetItems() [][]string {
	return i.items
}

type Model struct {
	table   *table.TableSingleType[string]
	infoBox *flexbox.FlexBox
	headers []string
}

func NewModel(data [][]string) *Model {
	m := Model{
		table:   table.NewTableSingleType[string](0, 0, headers),
		infoBox: flexbox.New(0, 0).SetHeight(7),
		headers: headers,
	}
	m.table.SetStylePassing(true)
	m.table.SetRatio(ratio).SetMinWidth(minSize)
	m.table.AddRows(data)

	infoText := help

	r1 := m.infoBox.NewRow()
	r1.AddCells(
		flexbox.NewCell(1, 1).
			SetID("info").
			SetContent(infoText),
		flexbox.NewCell(1, 1).
			SetID("info").
			SetContent(selectedValue).
			SetStyle(lipgloss.NewStyle().Bold(true)),
	)
	m.infoBox.AddRows([]*flexbox.Row{r1})
	m.infoBox.SetHeight(MIN_HEIGHT_HELP)
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

type Install struct {
	value string
}
type Update struct {
	value string
}
type Delete struct {
	value string
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Install:
		sendChange(m.infoBox.GetRow(0).GetCell(1), msg.value)
	case Update:
		sendChange(m.infoBox.GetRow(0).GetCell(1), msg.value)
	case Delete:
		sendChange(m.infoBox.GetRow(0).GetCell(1), msg.value)
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - m.infoBox.GetHeight())
		m.infoBox.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "down":
			m.table.CursorDown()
		case "up":
			m.table.CursorUp()
		case "left":
			m.table.CursorLeft()
		case "right":
			m.table.CursorRight()
		case "ctrl+s":
			x, _ := m.table.GetCursorLocation()
			m.table.OrderByColumn(x)
		case "i":
			name := m.table.GetValueOfColumnId(COLUMN_INDEX_TITLE)
			sendChange(m.infoBox.GetRow(0).GetCell(1), fmt.Sprintf("\n%s is installed", name))
			cmd := func() tea.Msg {
				clientMicrok8s.InstallModule(name)
				return Install{fmt.Sprintf("\n%s installed", name)}
			}
			return m, cmd
		case "d":
			name := m.table.GetValueOfColumnId(COLUMN_INDEX_TITLE)
			sendChange(m.infoBox.GetRow(0).GetCell(1), fmt.Sprintf("\n%s is deleted", name))
			cmd := func() tea.Msg {
				clientMicrok8s.RemoveModule(name)
				return Delete{fmt.Sprintf("\n%s deleted", name)}
			}
			return m, cmd
		case "u":
			name := m.table.GetValueOfColumnId(COLUMN_INDEX_TITLE)
			sendChange(m.infoBox.GetRow(0).GetCell(1), fmt.Sprintf("\n%s is updated", name))
			cmd := func() tea.Msg {
				return Update{fmt.Sprintf("\n%s updated", name)}
			}
			return m, cmd
		case "backspace":
			m.filterWithStr(msg.String())
		default:
			if len(msg.String()) == 1 {
				r := msg.Runes[0]
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					m.filterWithStr(msg.String())
				}
			}
		}
	}
	return m, nil
}

func (m *Model) filterWithStr(key string) {
	i, s := m.table.GetFilter()
	x, _ := m.table.GetCursorLocation()
	if x != i && key != "backspace" {
		m.table.SetFilter(x, key)
		return
	}
	if key == "backspace" {
		if len(s) == 1 {
			m.table.UnsetFilter()
			return
		} else if len(s) > 1 {
			s = s[0 : len(s)-1]
		} else {
			return
		}
	} else {
		s = s + key
	}
	m.table.SetFilter(i, s)
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.table.Render(), m.infoBox.Render())
}

func sendChange(cell *flexbox.Cell, text string) {
	cell.SetContent(text)
}
